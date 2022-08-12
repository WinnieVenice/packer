package cq

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/WinnieVenice/packer/conf"
	"github.com/WinnieVenice/packer/model"
	"github.com/WinnieVenice/packer/util"
)

var (
	server = NewServer()
)

func NewServer() *model.HttpServer {
	host := conf.V.GetString("service.cq.host")
	port := conf.V.GetInt("service.cq.port")
	return &model.HttpServer{
		Router: gin.Default(),
		Host:   host,
		Port:   port,
	}
}

func GetHostPort() string {
	return fmt.Sprintf("http://%s:%d", server.Host, server.Port)
}

func Run() {
	Register()
	err := server.Router.Run(fmt.Sprintf("%s:%d", server.Host, server.Port))
	if err != nil {
		fmt.Printf("Http Server run failed, err = (%s)\n", err.Error())
		panic(err)
	}
}

func Register() {
	server.Router.POST("/", DefaultHandle)
	server.Router.Static("pic", "./pic")
}

func DefaultHandle(c *gin.Context) {
	event := model.Event{}
	if err := c.ShouldBindJSON(&event); err != nil {
		fmt.Printf("http bind json err = (%+v)\n", err)
		return
	}
	// fmt.Println("get http req = ", event)

	msg := event.Message
	groupId := event.GroupId

	commCtx := map[string]string{
		"group_id": strconv.FormatInt(groupId, 10),
	}
	s := strings.Split(msg, " ")
	param := []string{}
	if len(s) > 1 {
		param = s[1:]
	}
	if f, ok := model.MsgHandlerMap[util.MatchCommand(s[0])]; ok {
		err := f.Handler(commCtx, param)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"status": "done"})
}
