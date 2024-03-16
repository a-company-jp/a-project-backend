package main

import (
	"a-project-backend/pkg/config"
	"a-project-backend/pkg/gcs"
	"a-project-backend/svc/pkg/handler"
	"a-project-backend/svc/pkg/middleware"
	"fmt"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	conf := config.Get()
	var dbUrl string
	switch conf.Infrastructure.MySQLDB.Protocol {
	case "tcp":
		dbUrl = fmt.Sprintf("%s:%s@%s(%s:%s)/%s?parseTime=true&loc=Asia%%2FTokyo&charset=utf8mb4&parseTime=True",
			conf.Infrastructure.MySQLDB.User,
			conf.Infrastructure.MySQLDB.Password,
			"tcp",
			conf.Infrastructure.MySQLDB.Host,
			conf.Infrastructure.MySQLDB.Port,
			conf.Infrastructure.MySQLDB.DBName)
	case "unix":
		dbUrl = fmt.Sprintf("%s:%s@%s(%s)/%s?parseTime=true&loc=Asia%%2FTokyo&charset=utf8mb4&parseTime=True",
			conf.Infrastructure.MySQLDB.User,
			conf.Infrastructure.MySQLDB.Password,
			"unix",
			conf.Infrastructure.MySQLDB.UnixSocket,
			conf.Infrastructure.MySQLDB.DBName)
	default:
		log.Fatalf("invalid protocol: %s", conf.Infrastructure.MySQLDB.Protocol)
	}

	db, err := gorm.Open(mysql.Open(dbUrl), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database, err: %v", err)
	}

	g, err := gcs.NewGCS()
	if err != nil {
		log.Fatalf("failed to connect to GCS, err: %v", err)
	}

	engine := gin.Default()
	engine.Use(cors.New(cors.Config{
		// アクセスを許可したいアクセス元
		AllowOrigins: []string{
			"http://localhost:3000",
		},
		// アクセスを許可したいHTTPメソッド(以下の例だとPUTやDELETEはアクセスできません)
		AllowMethods: []string{
			"POST",
			"GET",
			"PUT",
			"DELETE",
			"OPTIONS",
		},
		// 許可したいHTTPリクエストヘッダ
		AllowHeaders: []string{
			"Access-Control-Allow-Credentials",
			"Access-Control-Allow-Headers",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"Authorization",
		},
		// cookieなどの情報を必要とするかどうか
		AllowCredentials: true,
		// preflightリクエストの結果をキャッシュする時間
		MaxAge: 24 * time.Hour,
	}))

	engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	apiV1 := engine.Group("/api/v1")
	if err := Implement(apiV1, db, g); err != nil {
		log.Fatalf("Failed to start server... %v", err)
		return
	}

	if err := engine.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server... %v", err)
		return
	}
}

func Implement(rg *gin.RouterGroup, db *gorm.DB, g *gcs.GCS) error {
	//middlewareCROS := middleware.NewCORS()
	//middlewareCROS.ConfigureCORS(rg)

	middlewareAuth := middleware.NewAuth(db)
	authRg := rg.Use(middlewareAuth.VerifyUser())

	userHandler := handler.NewUser(db, g)
	milestoneHandler := handler.NewMileStone(db)

	authRg.Handle("GET", "/user/info/me", userHandler.GetMe())
	authRg.Handle("GET", "/user/info/:user_id", userHandler.GetUserInfo())
	authRg.Handle("GET", "/user/infos", userHandler.GetUserInfos())
	authRg.Handle("PUT", "/user/:user_id", userHandler.UpdateUserInfo())
	authRg.Handle("PUT", "/user/:user_id/icon", userHandler.UpdateUserIcon())
	authRg.Handle("POST", "/milestone", milestoneHandler.PostMileStone())
	authRg.Handle("PUT", "/milestone/:milestone_id", milestoneHandler.UpdateMileStone())
	authRg.Handle("DELETE", "/milestone/:milestone_id", milestoneHandler.DeleteMileStone())

	return nil
}
