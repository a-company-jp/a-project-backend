package main

import (
	"a-project-backend/pkg/config"
	"a-project-backend/svc/pkg/handler"
	"a-project-backend/svc/pkg/middleware"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
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

	engine := gin.Default()
	apiV1 := engine.Group("/api/v1")
	if err := Implement(apiV1, db); err != nil {
		log.Fatalf("Failed to start server... %v", err)
		return
	}

	if err := engine.Run("localhost:8080"); err != nil {
		log.Fatalf("Failed to start server... %v", err)
		return
	}
}

func Implement(rg *gin.RouterGroup, db *gorm.DB) error {
	middlewareCROS := middleware.NewCORS()
	middlewareCROS.ConfigureCORS(rg)

	middlewareAuth := middleware.NewAuth()
	authRg := rg.Use(middlewareAuth.VerifyUser())

	userHandler := handler.NewUser(db)
	authRg.Handle("GET", "/user/info", userHandler.GetUserInfo())
	return nil
}
