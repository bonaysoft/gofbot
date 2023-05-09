package main

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

type Map map[string]interface{}

type Server struct {
	*gin.Engine
}

func New(robots []*Robot) (*Server, error) {
	buildHandler := func(robot Robot) gin.HandlerFunc {
		return func(c *gin.Context) {
			incomingHandler(c, &robot)
		}
	}

	router := gin.Default()
	for _, robot := range robots {
		router.POST(fmt.Sprintf("/incoming/%s", robot.Alias), buildHandler(*robot))
	}

	return &Server{
		Engine: router,
	}, nil
}

func incomingHandler(ctx *gin.Context, robot *Robot) {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	params := make(Map)
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
		if resp, err := http.DefaultClient.Post(robot.WebHook, "application/json", body); err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		} else {
			defer resp.Body.Close()
			rb, _ := ioutil.ReadAll(resp.Body)
			ctx.Data(resp.StatusCode, resp.Header.Get("Content-Type"), rb)
		}
	}
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

		if nextParams, ok := params[k].(map[string]interface{}); ok {
			params = nextParams
		}
	}

	return ""
}
