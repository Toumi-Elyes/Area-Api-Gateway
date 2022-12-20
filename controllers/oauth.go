package controllers

import (
	"api-gateway/models"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

type AccessTokenResponse struct {
	Success bool
	Url     string
}

// create a route GET /Oauth/:service-name/ to get the url to redirect the user to the oauth provider
func GetOauthLink(c *gin.Context) {
	serviceName := c.Param("service-name")

	if serviceName == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	// call the dispatcher to get the url
	dispatcherUrl := os.Getenv("DISPATCHER_URL")
	if dispatcherUrl == "" {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	res, _ := http.Get(dispatcherUrl + "/v1/oauth/" + serviceName)
	byte_body, _ := ioutil.ReadAll(res.Body)
	var data AccessTokenResponse
	json.Unmarshal(byte_body, &data)
	c.JSON(http.StatusOK, data)
}

type tokenRequest struct {
	Code   string `json:"code"`
	UserId string `json:"user_id"`
}

func StoreToken(c *gin.Context) {
	serviceName := c.Param("service-name")
	token_id, _, _ := ExtractToken(c)
	user_id, _ := uuid.Parse(token_id)
	token := tokenRequest{UserId: user_id.String()}
	err := c.ShouldBindBodyWith(&token, binding.JSON)

	jsonData, err := json.Marshal(token)
	if serviceName == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	// call the dispatcher to get the url
	dispatcherUrl := os.Getenv("DISPATCHER_URL")
	if dispatcherUrl == "" {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	resp, err := http.Post(dispatcherUrl+"/v1/oauth/"+serviceName, "application/json", bytes.NewBuffer(jsonData))
	if err != nil || resp.StatusCode != http.StatusOK {
		// logger.Log.Errorf("Failed Create token from your code.\nError message: [%v]\n", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	service, _ := models.GetServicesByName(serviceName)
	models.AddServiceToUser(models.UserModel{Id: user_id}, service)
	c.JSON(http.StatusOK, "{\"message\": \"Oauth setup successfully\"}")
}

func RemoveAuth(c *gin.Context) {
	serviceName := c.Param("service-name")
	token_id, _, _ := ExtractToken(c)
	user_id, _ := uuid.Parse(token_id)

	if serviceName == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	service, _ := models.GetServicesByName(serviceName)
	models.RemoveServiceFromUser(models.UserModel{Id: user_id}, service)
	c.JSON(http.StatusOK, "{\"message\": \"Oauth remove successfully\"}")
}
