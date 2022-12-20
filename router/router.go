// Package router package store all the routes of the api-gateway
package router

import (
	"api-gateway/controllers"
	"api-gateway/middlewares"

	"github.com/gin-gonic/gin"
)

// Apply the routes to the corresponding controllers
func Apply(r *gin.Engine) {
	r.GET("/ping", controllers.Ping)

	// About.json
	r.GET("/about.json", controllers.AboutJson)

	// Users
	r.POST("/register", controllers.Register)
	r.GET("/readUsers", middlewares.CheckAdmin, controllers.ReadUsers)
	r.GET("/readUser", middlewares.CheckAdmin, controllers.ReadUser)
	r.POST("/updateUser", middlewares.CheckAdmin, controllers.UpdateUser)
	r.POST("/deleteUser", middlewares.CheckAdmin, controllers.DeleteUser)

	// Login - Session
	r.POST("/login", controllers.Login)
	r.POST("/loginAdmin", controllers.LoginAdmin)
	r.GET("/session", middlewares.CheckSession, controllers.Session)

	//only used by the dispatcher
	r.POST("/services", controllers.CreateServices)

	// Services
	r.DELETE("/services", controllers.DeleteServices)

	// Areas
	r.POST("/area", middlewares.CheckSession, controllers.CreateArea)
	r.GET("/area", middlewares.CheckSession, controllers.GetUserAreas)
	r.POST("/area/delete", middlewares.CheckSession, controllers.DeleteArea)
	r.POST("/areas/delete", middlewares.CheckAdmin, controllers.DeleteUserAreas)

	// Oauth
	r.GET("/oauth/:service-name/", middlewares.CheckSession, controllers.GetOauthLink)
	r.POST("/oauth/:service-name/", middlewares.CheckSession, controllers.StoreToken)
	r.DELETE("/oauth/:service-name/", middlewares.CheckSession, controllers.RemoveAuth)

	//user
	r.GET("/user/services", middlewares.CheckSession, controllers.UserServicesList)

	// Get all Service information for the client
	r.GET("/about", controllers.GetAllAboutInfo)

	r.NoRoute(controllers.NotFound)
}
