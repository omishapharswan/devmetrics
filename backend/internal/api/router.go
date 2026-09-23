package api

import (
	"github.com/gin-gonic/gin"

	"github.com/omishapharswan/devmetrics/backend/internal/api/handlers"
)

// NewRouter builds the Gin engine with all API routes registered.
func NewRouter() *gin.Engine {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	api := r.Group("/api")
	{
		api.GET("/health", handlers.Health)
	}

	return r
}
