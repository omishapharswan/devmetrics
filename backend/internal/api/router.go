package api

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"github.com/omishapharswan/devmetrics/backend/internal/api/handlers"
)

// dbContextKey is the gin.Context key under which the shared *sql.DB
// connection is stored for handlers to retrieve.
const dbContextKey = "db"

// NewRouter builds the Gin engine with all API routes registered.
func NewRouter(conn *sql.DB) *gin.Engine {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Set(dbContextKey, conn)
		c.Next()
	})

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
