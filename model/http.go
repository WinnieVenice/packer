package model

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	Router *gin.Engine
	Ip     string
	Port   string
	Url    string
}

type HttpClient struct {
	Client http.Client
}

type Handler struct{}

type Event struct {
	Message string `json:"message"`
	GroupId int64  `json:"group_id"`
}
