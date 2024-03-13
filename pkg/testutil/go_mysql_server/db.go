package go_mysql_server

import (
	"context"
	"fmt"
	"github.com/dolthub/vitess/go/mysql"

	sqle "github.com/dolthub/go-mysql-server"
	"github.com/dolthub/go-mysql-server/memory"
	"github.com/dolthub/go-mysql-server/server"
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/mysql_db"
)

type GoMysqlServer struct {
	Address string
	Port    int
	DBName  string
}

func (g GoMysqlServer) Run(
	ctx context.Context,
) {
	db := memory.NewDatabase(g.DBName)
	db.BaseDatabase.EnablePrimaryKeyIndexes()

	pro := memory.NewDBProvider(db)
	session := memory.NewSession(sql.NewBaseSession(), pro)
	sqlCtx := sql.NewContext(ctx, sql.WithSession(session))
	engine := sqle.NewDefault(pro)

	sqlCtx.SetCurrentDatabase(g.DBName)
	config := server.Config{
		Protocol: "tcp",
		Address:  fmt.Sprintf("%s:%d", g.Address, g.Port),
	}

	s, err := server.NewServer(config, engine, sessionBuilder(pro), nil)
	if err != nil {
		panic(err)
	}

	go func() {
		if err = s.Start(); err != nil {
			panic(err)
		}
	}()
}

func sessionBuilder(pro *memory.DbProvider) server.SessionBuilder {
	return func(ctx context.Context, conn *mysql.Conn, addr string) (sql.Session, error) {
		host := ""
		user := ""
		mysqlConnectionUser, ok := conn.UserData.(mysql_db.MysqlConnectionUser)
		if ok {
			host = mysqlConnectionUser.Host
			user = mysqlConnectionUser.User
		}
		client := sql.Client{Address: host, User: user, Capabilities: conn.Capabilities}
		baseSession := sql.NewBaseSessionWithClientServer(addr, client, conn.ConnectionID)
		return memory.NewSession(baseSession, pro), nil
	}
}
