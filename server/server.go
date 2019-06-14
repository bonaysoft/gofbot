package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var tplArgExp = regexp.MustCompile(`{{(\s*\$\S+\s*)}}`)

type Server struct {
	router *gin.Engine
	robots []*Robot
}

func New() (*Server, error) {
	robot, err := newRobot("conf/conf.yaml")
	if err != nil {
		return nil, err
	}

	return &Server{
		router: gin.Default(),
		robots: []*Robot{robot},
	}, nil
}

func (s *Server) Run() error {
	for _, robot := range s.robots {
		s.router.POST(fmt.Sprintf("/incoming/%s", robot.Alias), func(context *gin.Context) {
			s.incomingHandler(context, robot)
		})
	}

	return s.router.Run(":9613")
}

func (s *Server) incomingHandler(ctx *gin.Context, robot *Robot) {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	params := make(map[string]interface{})
	if err := json.Unmarshal(body, &params); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	for _, msg := range robot.Messages {
		if !msg.Exp.Match(body) {
			continue
		}

		// 正则替换参数
		message := buildMessage(msg.Template, params)
		body := buildPostBody(robot.BodyTpl, message)
		http.DefaultClient.Timeout = 3 * time.Second
		if _, err := http.DefaultClient.Post(robot.WebHook, "application/json", body); err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}
}

type variable struct {
	full string
	name string
}

func buildMessage(tpl string, params map[string]interface{}) string {
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

func extractArgs(params map[string]interface{}, key string) string {
	key = strings.Replace(key, "$", "", -1)
	keys := strings.Split(key, ".")
	for index, k := range keys {
		if index == len(keys)-1 {
			return params[k].(string)
		}

		if nextParams, ok := params[k].(map[string]interface{}); ok {
			params = nextParams
		}
	}

	return ""
}
