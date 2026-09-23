package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Health reports that the API server and its database connection are up.
func Health(c *gin.Context) {
	dbStatus := "ok"
	if err := DBFromContext(c).Ping(); err != nil {
		dbStatus = "unreachable"
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "devmetrics-backend",
		"db":      dbStatus,
	})
}
