package controllers

import (
	"api-gateway/logger"
	"api-gateway/models"
	"net/http"
	"github.com/google/uuid"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

var Services []models.ServiceModel

func CreateServices(c *gin.Context) {
	aboutData := models.AboutModel{}
	err := c.ShouldBindBodyWith(&aboutData, binding.JSON)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	Services = aboutData.Server.Services

	for _, service := range aboutData.Server.Services {
		_, err = models.CreateService(service)
		if err != nil {
			logger.Log.Errorf("Failed to create the services.\nError message: [%v]\n", err)
			c.AbortWithStatus(http.StatusNotAcceptable)
			return
		}
	}
	c.String(http.StatusOK, "Services created successfully")
}

func DeleteServices(c *gin.Context) {
	_, err := models.DeleteServices()

	if err != nil {
		logger.Log.Errorf("Failed to delete.\nError message: [%v]\n", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	Services = []models.ServiceModel{}
	c.String(http.StatusOK, "Services deleted successfully")
}


func UserServicesList(c *gin.Context) {
	token_id, _, _ := ExtractToken(c)
	tmp_id, _ := uuid.Parse(token_id)
	user := models.UserModel{Id: tmp_id}
	userDb, _ := models.ReadOneUser(user)
	userServiceDb, _ := models.GetUserServices(userDb)

	var userServiceList = models.UserServicesRequest{Services: []string{}}
	for _, service := range userServiceDb {

		userServiceList.Services = append(userServiceList.Services, service.Name)
	}
	c.JSON(http.StatusOK,  userServiceList)
}
