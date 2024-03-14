package main

import (
	"a-project-backend/pkg/config"
	"a-project-backend/pkg/gcs"
	"a-project-backend/svc/pkg/handler"
	"a-project-backend/svc/pkg/middleware"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	conf := config.Get()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.Infrastructure.MySQLDB.User,
		conf.Infrastructure.MySQLDB.Password,
		conf.Infrastructure.MySQLDB.Host,
		conf.Infrastructure.MySQLDB.Port,
		conf.Infrastructure.MySQLDB.DBName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database, err: %v", err)
	}
	g, err := gcs.NewGCS()
	if err != nil {
		log.Fatalf("failed to connect to GCS, err: %v", err)
	}

	engine := gin.Default()
	apiV1 := engine.Group("/api/v1")
	if err := Implement(apiV1, db, g); err != nil {
		log.Fatalf("Failed to start server... %v", err)
		return
	}

	if err := engine.Run("localhost:8080"); err != nil {
		log.Fatalf("Failed to start server... %v", err)
		return
	}
}

func Implement(rg *gin.RouterGroup, db *gorm.DB, g *gcs.GCS) error {
	middlewareCROS := middleware.NewCORS()
	middlewareCROS.ConfigureCORS(rg)

	middlewareAuth := middleware.NewAuth(db)
	authRg := rg.Use(middlewareAuth.VerifyUser())

	userHandler := handler.NewUser(db, g)
	authRg.Handle("GET", "/user/info", userHandler.GetUserInfo())
	return nil
}
