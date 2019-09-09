package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/saltbo/gofbot/robot"
)

type Server struct {
	http.Server
	router *gin.Engine
	robots map[string]*robot.Robot
}

func NewServer() *Server {
	router := gin.Default()
	server := &Server{
		Server: http.Server{
			Addr:    ":8080",
			Handler: router,
		},
		router: router,
		robots: make(map[string]*robot.Robot),
	}

	router.POST("/incoming/:alias", server.incomingHandler)
	return server
}

func (s *Server) SetupRobots(robots []*robot.Robot) {
	for _, bot := range robots {
		s.robots[bot.Alias] = bot
	}
}

func (s *Server) Run(addr ...string) error {
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.Server.Shutdown(ctx)
}

func (s *Server) incomingHandler(ctx *gin.Context) {
	alias := ctx.Param("alias")
	bot, ok := s.robots[alias]
	if !ok {
		ctx.AbortWithError(http.StatusNotFound, fmt.Errorf("not found your robot"))
		return
	}

	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	params := make(robot.Map)
	if err := json.Unmarshal(body, &params); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	msg, err := bot.MatchMessage(body)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	msgStr := robot.BuildMessage(msg.Template, params) // 正则替换参数
	postBody := robot.BuildPostBody(bot.BodyTpl, msgStr)
	if err := forwardToRobot(ctx, bot.WebHook, postBody); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}

func forwardToRobot(ctx *gin.Context, url string, body io.Reader) error {
	http.DefaultClient.Timeout = 3 * time.Second
	resp, err := http.DefaultClient.Post(url, "application/json", body)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	rb, _ := ioutil.ReadAll(resp.Body)
	ctx.Data(resp.StatusCode, resp.Header.Get("Content-Type"), rb)
	return nil
}
