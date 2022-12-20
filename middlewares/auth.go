package middlewares

import (
	"api-gateway/controllers"
	"api-gateway/logger"
	"api-gateway/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckSession(c *gin.Context) {
	tokenId, tokenEmail, err := controllers.ExtractToken(c)
	log.Println(tokenEmail, tokenId)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	users, err := models.ReadUsers()

	if err != nil {
		logger.Log.Errorf("Failed to read the User.\nError message: [%v]\n", err)
		c.AbortWithStatus(http.StatusNotAcceptable)
		return
	}

	for _, u := range users {
		if err != nil {
			c.AbortWithStatus(http.StatusNotAcceptable)
			return
		}
		if u.Email == tokenEmail && u.ID.String() == tokenId {
			c.Next()
			return
		}
	}
	c.AbortWithStatus(http.StatusNotAcceptable)
}
