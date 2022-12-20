// Package controllers provide controllers for the backend routes
package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// NotFound Return the http status not found (400) with the json: {"message": "Route not found"}
func NotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"message": "Route Not Found",
	})
}
