package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/omishapharswan/devmetrics/backend/internal/service"
)

// ListScans returns every persisted scan, most recent first.
func ListScans(c *gin.Context) {
	scans, err := service.ListScans(DBFromContext(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"scans": scans})
}

// GetScan returns the full result (scan, files, dependencies) for one
// scan ID.
func GetScan(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid scan id"})
		return
	}

	result, err := service.GetScan(DBFromContext(c), id)
	if err != nil {
		if errors.Is(err, service.ErrScanNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
