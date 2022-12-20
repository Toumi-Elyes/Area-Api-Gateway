// Package controllers provide controllers for the backend routes
package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Ping Route: GET "/ping"
// Send a string that contains "Pong" and the http status ok (200). No parameters are needed.
func Ping(c *gin.Context) {
	c.String(http.StatusOK, "Pong")
}
