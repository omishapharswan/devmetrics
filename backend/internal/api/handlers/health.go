package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Health reports that the API server is up and reachable.
func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "devmetrics-backend",
	})
}
