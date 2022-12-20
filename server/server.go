// Package server package handle the server creation and maintenance
package server

import (
	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
)

type Server struct {
	Def *gin.Engine
}

func NewServer() Server {
	router := gin.Default()

	router.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		Credentials:     false,
		ValidateHeaders: false,
	}))
	server := Server{Def: router}
	return server
}
