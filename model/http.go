package model

import (
	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	Router *gin.Engine
	Host   string
	Port   int
}

type Handler struct{}

type Event struct {
	Message string `json:"message"`
	GroupId int64  `json:"group_id"`
}
