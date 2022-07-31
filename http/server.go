package http

import (
	"fmt"
	"net/http"
	"packer/model"
	"packer/util"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	server *model.HttpServer
)

func init() {
	server = NewServer("", "")
}

func NewServer(ip string, port string) *model.HttpServer {
	if len(ip) <= 0 {
		ip = "localhost"
	}
	if len(port) <= 0 {
		port = "5701"
	}

	server := model.HttpServer{
		Router: gin.Default(),
		Ip:     ip,
		Port:   port,
		Url:    fmt.Sprintf("http://%s:%s", ip, port),
	}
	return &server
}

func GetServer() *model.HttpServer {
	if server == nil {
		server = NewServer("", "")
	}
	return server
}

func Run() error {
	Register()
	err := server.Router.Run(fmt.Sprintf(":%s", server.Port))
	if err != nil {
		fmt.Printf("Http Server run failed, err = (%s)\n", err.Error())
		return err
	}
	return nil
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
		err := f(commCtx, param)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"status": "done"})
}
