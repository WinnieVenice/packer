package model

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

const (
	CQ_AT = "[CQ:at,qq=%d]"
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
	UserId  int64  `json:"user_id"`
}

func CqCodeAt(id int64) string {
	return fmt.Sprintf(CQ_AT, id)
}
