package controllers

import (
	"api-gateway/logger"
	"api-gateway/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
)

func CreateArea(c *gin.Context) {
	areaData := models.AreaRequestModel{}
	err := c.ShouldBindBodyWith(&areaData, binding.JSON)
	if err != nil {
		logger.Log.Errorf("Failed to bind the area data.\nError message: [%v]\n", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	token_id, _, err := ExtractToken(c)

	if err != nil {
		c.String(http.StatusUnauthorized, "Unauthorized")
		return
	}
	tmp_id, err := uuid.Parse(token_id)

	if err != nil {
		logger.Log.Errorf("Failed to parse the token id to uuid.\nError message: [%v]\n", err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	areaData.Payload.UserId = tmp_id

	dataOfTheModel, err := models.CreateArea(areaData)
	if err != nil {
		logger.Log.Errorf("Failed to create the area.\nError message: [%v]\n", err)
		c.AbortWithStatus(http.StatusNotAcceptable)
		return
	}
	areaData.Payload.AreaId = dataOfTheModel.AreaID
	jsonData, err := json.Marshal(areaData)
	if err != nil {
		logger.Log.Errorf("Failed to convert the data to json.\nError message: [%v]\n", err)
		c.AbortWithStatus(http.StatusNotAcceptable)
		return
	}
	dispatcherUrl := os.Getenv("DISPATCHER_URL")
	if dispatcherUrl == "" {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	resp, err := http.Post(dispatcherUrl+"/v1/area", "application/json", bytes.NewBuffer(jsonData))
	var data models.DataResponse
	byte_body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(byte_body, &data)

	if err != nil || data.Success == false {
		if strings.Contains(data.Payload, "Bad creds") {
			token_id, _, _ := ExtractToken(c)
			user_id, _ := uuid.Parse(token_id)
			service, _ := models.GetServicesByName(strings.Split(data.Payload, ":")[1])
			models.RemoveServiceFromUser(models.UserModel{Id: user_id}, service)
		}
		logger.Log.Errorf("Failed to send the area to the dispatcher.\nError message: [%v]\n", err)
		c.AbortWithStatus(http.StatusNotAcceptable)
		return
	}
	// check status code
	if resp.StatusCode != http.StatusOK {
		logger.Log.Errorf("Failed to send the area to the dispatcher.\nError message: [%v]\n", err)
		c.AbortWithStatus(http.StatusNotAcceptable)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Log.Errorf("Failed to close the response body.\nError message: [%v]\n", err)
		}
	}(resp.Body)
	c.JSON(http.StatusOK, "{\"message\": \"Area created successfully\"}")
}

func GetUserAreas(c *gin.Context) {
	userData := models.UserModel{}
	c.ShouldBindBodyWith(&userData, binding.JSON)

	tokenId, _, err := ExtractToken(c)

	if err != nil {
		c.String(http.StatusUnauthorized, "Unauthorized")
		return
	}
	tmpId, err := uuid.Parse(tokenId)

	if err != nil {
		logger.Log.Errorf("Failed to parse the token id to uuid.\nError message: [%v]\n", err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	areas, err := models.GetUserAreas(models.UserModel{Id: tmpId})

	if areas != nil && len(areas) > 0 && !models.CheckAreas(tmpId, areas) {
		logger.Log.Errorf("Failed to get the user areas.\nError message: [%v]\n", err)
		c.AbortWithStatus(http.StatusNotAcceptable)
	}
	if err != nil {
		logger.Log.Errorf("Failed to get the user areas.\nError message: [%v]\n", err)
		c.AbortWithStatus(http.StatusNotAcceptable)
	} else {
		userAreas := models.UserAreasProps{Areas: []models.UserArea{}}
		for _, area := range areas {
			var ActionReactionList []models.ArOrRea
			var areaList []models.ActionReactionModel
			json.Unmarshal([]byte(area.ActionReaction), &areaList)

			for _, ArOrRea := range areaList {
				var fieldList []models.ArOrReaField
				for key, v := range ArOrRea.Payload {
					fieldList = append(fieldList, models.ArOrReaField{Name: key, Value: v})
				}
				dispatcherUrl := os.Getenv("DISPATCHER_URL")
				if dispatcherUrl == "" {
					c.AbortWithStatus(http.StatusInternalServerError)
					return
				}
				iconResponse, err := http.Get(dispatcherUrl + "/v1/icon/" + ArOrRea.Service)
				var icon string = ""

				if err == nil {
					var data models.IconResponse
					byte_body, _ := io.ReadAll(iconResponse.Body)
					json.Unmarshal(byte_body, &data)
					icon = data.Icon
				}
				ActionReactionList = append(ActionReactionList, models.ArOrRea{Type: ArOrRea.Type, Service: ArOrRea.Service, Name: ArOrRea.Name, Order: ArOrRea.Order, Icon: icon, Fields: fieldList})
			}
			userAreas.Areas = append(userAreas.Areas, models.UserArea{AreaId: area.AreaID, AreaName: area.AreaName, ActionReaction: ActionReactionList})
		}
		c.JSON(http.StatusOK, userAreas)
	}
}

func DeleteArea(c *gin.Context) {
	areaData := models.DeleteAreaRequest{}
	err := c.ShouldBindBodyWith(&areaData, binding.JSON)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}
	areas, err := models.GetAreasById(areaData.AreaId)

	nbRow, err := models.DeleteArea(areaData.AreaId)
	if err != nil {
		logger.Log.Errorf("Failed to delete area.\nError message: [%v]\n", err)
		c.AbortWithStatus(http.StatusNotAcceptable)
	} else if nbRow == 0 {
		c.String(http.StatusBadRequest, fmt.Sprintf("Doesn't find any Area with this id."))
		return
	}

	var nextDelete models.AreaDeleteRequestDispatcher
	var areaList []models.ActionReactionModel
	json.Unmarshal([]byte(areas.ActionReaction), &areaList)
	for _, ArOrRea := range areaList {
		nextDelete.AreaList = append(nextDelete.AreaList, models.AreaToDelete{Type: ArOrRea.Type, ServiceName: ArOrRea.Service})
	}

	nextDelete.AreaId = areaData.AreaId
	jsonData, err := json.Marshal(nextDelete)
	if err != nil {
		logger.Log.Errorf("Failed to convert the data to json.\nError message: [%v]\n", err)
		c.AbortWithStatus(http.StatusNotAcceptable)
		return
	}
	fmt.Println(string(jsonData))
	dispatcherUrl := os.Getenv("DISPATCHER_URL")
	if dispatcherUrl == "" {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	req, err := http.NewRequest("POST", dispatcherUrl+"/v1/area/delete", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println(err)
		return
	}
	client := &http.Client{}
	client.Do(req)
	c.JSON(http.StatusOK, "{\"message\": \"Area deleted successfully\"}")
}

func DeleteUserAreas(c *gin.Context) {
	userData := models.UserModel{}
	err := c.ShouldBindBodyWith(&userData, binding.JSON)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	usersAreas, err := models.GetUserAreas(userData)

	for _, areaData := range usersAreas {
		areas, err := models.GetAreasById(areaData.AreaID)

		nbRow, err := models.DeleteArea(areaData.AreaID)
		if err != nil {
			logger.Log.Errorf("Failed to delete area.\nError message: [%v]\n", err)
			c.AbortWithStatus(http.StatusNotAcceptable)
		} else if nbRow == 0 {
			c.String(http.StatusBadRequest, fmt.Sprintf("Doesn't find any Area with this id."))
			return
		}

		var nextDelete models.AreaDeleteRequestDispatcher
		var areaList []models.ActionReactionModel
		json.Unmarshal([]byte(areas.ActionReaction), &areaList)
		for _, ArOrRea := range areaList {
			nextDelete.AreaList = append(nextDelete.AreaList, models.AreaToDelete{Type: ArOrRea.Type, ServiceName: ArOrRea.Service})
		}

		nextDelete.AreaId = areaData.AreaID
		jsonData, err := json.Marshal(nextDelete)
		if err != nil {
			logger.Log.Errorf("Failed to convert the data to json.\nError message: [%v]\n", err)
			c.AbortWithStatus(http.StatusNotAcceptable)
			return
		}
		fmt.Println(string(jsonData))
		dispatcherUrl := os.Getenv("DISPATCHER_URL")
		if dispatcherUrl == "" {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		req, err := http.NewRequest("POST", dispatcherUrl+"/v1/area/delete", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Println(err)
			return
		}
		client := &http.Client{}
		client.Do(req)
	}
	c.JSON(http.StatusOK, "{\"message\": \"User's Areas deleted successfully\"}")
}
