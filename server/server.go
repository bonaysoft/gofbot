package main

import (
	"bytes"
	`context`
	"encoding/json"
	`fmt`
	`io`
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var tplArgExp = regexp.MustCompile(`{{(\s*\$\S+\s*)}}`)

type Map map[string]interface{}

type Server struct {
	http.Server
	robots map[string]*Robot
}

func NewServer() *Server {
	router := gin.Default()
	server := &Server{
		Server: http.Server{
			Addr:    ":8080",
			Handler: router,
		},
		robots: make(map[string]*Robot),
	}

	router.POST("/incoming/:alias", server.incomingHandler)
	return server
}

func (s *Server) SetupRobots(robots []*Robot) {
	for _, robot := range robots {
		s.robots[robot.Alias] = robot
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
	robot, ok := s.robots[alias]
	if !ok {
		ctx.AbortWithError(http.StatusNotFound, fmt.Errorf("not found your robot"))
		return
	}

	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	params := make(Map)
	if err := json.Unmarshal(body, &params); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	msg, err := robot.MatchMessage(body)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	msgStr := buildMessage(msg.Template, params) // 正则替换参数
	postBody := buildPostBody(robot.BodyTpl, msgStr)
	if err := forwardToRobot(ctx, robot.WebHook, postBody); err != nil {
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

type variable struct {
	full string
	name string
}

func buildMessage(tpl string, params Map) string {
	variables := make([]variable, 0)
	for _, v := range tplArgExp.FindAllStringSubmatch(tpl, -1) {
		variables = append(variables, variable{full: v[0], name: v[1]})
	}

	newMsg := tpl
	for _, v := range variables {
		newMsg = strings.Replace(newMsg, v.full, extractArgs(params, strings.TrimSpace(v.name)), -1)
	}

	if strconv.CanBackquote(newMsg) {
		return newMsg
	}

	return strconv.Quote(newMsg)
}

func buildPostBody(bodyTpl string, message string) *bytes.Buffer {
	return bytes.NewBufferString(strings.Replace(bodyTpl, "$template", message, -1))
}

func extractArgs(params Map, key string) string {
	key = strings.Replace(key, "$", "", -1)
	keys := strings.Split(key, ".")
	for index, k := range keys {
		if index == len(keys)-1 {
			if v, ok := params[k].(string); ok {
				return v
			}
			return ""
		}

		if nextParams, ok := params[k].(Map); ok {
			params = nextParams
		}
	}

	return ""
}
