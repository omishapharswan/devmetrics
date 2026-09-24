package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/omishapharswan/devmetrics/backend/internal/service"
)

type scanRequest struct {
	FolderPath string `json:"folderPath" binding:"required"`
}

// Scan runs a full DevMetrics scan of the requested folder, persists the
// result, and returns the scan summary and per-file metrics.
func Scan(c *gin.Context) {
	var req scanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "folderPath is required"})
		return
	}

	result, err := service.RunScan(req.FolderPath, DBFromContext(c))
	if err != nil {
		if errors.Is(err, service.ErrFolderNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}
