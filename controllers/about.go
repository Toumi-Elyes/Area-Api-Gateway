package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func GetAllAboutInfo(c *gin.Context) {
	// call the dispatcher to get the url
	dispatcherUrl := os.Getenv("DISPATCHER_URL")
	if dispatcherUrl == "" {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	res, _ := http.Get(dispatcherUrl + "/v1/area/info")

	var jsonMap map[string]interface{}
	byte_body, _ := io.ReadAll(res.Body)
	err := json.Unmarshal(byte_body, &jsonMap)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, jsonMap)
}

type OldTypeInfo struct {
	Services []Service `json:"services,omitempty"`
}

type Service struct {
	Name                string   `json:"name,omitempty"`
	Label               string   `json:"label,omitempty"`
	Icon                string   `json:"icon,omitempty"`
	Actions             []Action `json:"actions,omitempty"`
	Reactions           []Action `json:"reactions,omitempty"`
	NeedsAuthentication bool     `json:"needsAuthentication,omitempty"`
}

type Action struct {
	Name        string  `json:"name,omitempty"`
	Label       string  `json:"label,omitempty"`
	Description string  `json:"description,omitempty"`
	Fields      []Field `json:"fields,omitempty"`
}

type Field struct {
	Name    string `json:"name,omitempty"`
	Label   string `json:"label,omitempty"`
	Type    Type   `json:"type,omitempty"`
	Tooltip string `json:"tooltip,omitempty"`
}

type Type string

const (
	Number Type = "number"
	String Type = "string"
)

type NewTypeInfo struct {
	Client Client `json:"client,omitempty"`
	Server Server `json:"server,omitempty"`
}

type Client struct {
	Host string `json:"host,omitempty"`
}

type Server struct {
	CurrentTime int64         `json:"current_time,omitempty"`
	Services    []ServicesNew `json:"services,omitempty"`
}

type ServicesNew struct {
	Name      string    `json:"name,omitempty"`
	Actions   []Actions `json:"actions"`
	Reactions []Actions `json:"reactions"`
}

type Actions struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

func convertToNewType(oldType OldTypeInfo, ipOfClient string) NewTypeInfo {
	var newType NewTypeInfo
	newType.Client.Host = ipOfClient
	newType.Server.CurrentTime = CurrentTime()
	for _, service := range oldType.Services {
		var newService ServicesNew
		newService.Name = service.Name
		for _, action := range service.Actions {
			var newAction Actions
			newAction.Name = action.Name
			newAction.Description = action.Description
			newService.Actions = append(newService.Actions, newAction)
		}
		for _, reaction := range service.Reactions {
			var newReaction Actions
			newReaction.Name = reaction.Name
			newReaction.Description = reaction.Description
			newService.Reactions = append(newService.Reactions, newReaction)
		}
		newType.Server.Services = append(newType.Server.Services, newService)
	}
	for i := range newType.Server.Services {
		if newType.Server.Services[i].Actions == nil {
			newType.Server.Services[i].Actions = []Actions{}
		}
		if newType.Server.Services[i].Reactions == nil {
			newType.Server.Services[i].Reactions = []Actions{}
		}
	}
	return newType
}

func CurrentTime() int64 {
	return time.Now().UnixNano() / int64(time.Second)
}

func AboutJson(c *gin.Context) {
	dispatcherUrl := os.Getenv("DISPATCHER_URL")
	if dispatcherUrl == "" {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	res, _ := http.Get(dispatcherUrl + "/v1/area/info")

	var jsonMap map[string]interface{}
	byteBody, _ := io.ReadAll(res.Body)
	err := json.Unmarshal(byteBody, &jsonMap)
	if err != nil {
		return
	}
	fmt.Println(jsonMap)
	var oldType OldTypeInfo
	err = json.Unmarshal(byteBody, &oldType)
	if err != nil {
		return
	}
	ipOfClient := c.ClientIP()
	fmt.Println(ipOfClient)
	newType := convertToNewType(oldType, ipOfClient)
	c.JSON(http.StatusOK, newType)
}
