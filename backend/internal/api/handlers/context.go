package handlers

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

// DBFromContext retrieves the shared *sql.DB connection stored on the
// gin.Context by the router's middleware.
func DBFromContext(c *gin.Context) *sql.DB {
	return c.MustGet("db").(*sql.DB)
}
