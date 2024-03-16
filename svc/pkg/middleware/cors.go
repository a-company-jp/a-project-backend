package middleware

import (
	"github.com/gin-gonic/gin"
)

type CORS struct {
	targetHost string
}

func NewCORS() CORS {
	return CORS{targetHost: "https://main.a.shion.pro"}
}

func (cr CORS) ConfigureCORS(rg *gin.RouterGroup) {
	rg.Use(cr.middleware())
	// this does absolutely nothing because OPTIONS request will be intercepted by the middleware,
	// but this is needed to listen for OPTIONS requests
	rg.OPTIONS("/*path", func(c *gin.Context) {
		c.AbortWithStatus(200)
	})
}

func (cr CORS) middleware() gin.HandlerFunc {
	allowedOrigins := []string{cr.targetHost, "http://localhost:3000"}
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin == "" {
			origin = c.Request.Header.Get("Referer")
		}
		allowedOrigin := ""
		for _, o := range allowedOrigins {
			if origin == o || origin == o+"/" {
				allowedOrigin = origin
				break
			}
		}
		if allowedOrigin == "" {
			allowedOrigin = allowedOrigins[0]
		}
		c.Header("Access-Control-Allow-Origin", allowedOrigin)
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}
		c.Next()
	}
}
