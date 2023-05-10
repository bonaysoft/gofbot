package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bonaysoft/gofbot/pkg/robot"
	"github.com/gin-gonic/gin"
)

type Server struct {
	*gin.Engine
}

func New(robots []*robot.Robot) (*Server, error) {
	buildHandler := func(robot robot.Robot) gin.HandlerFunc {
		return func(c *gin.Context) {
			incomingHandler(c, &robot)
		}
	}

	router := gin.Default()
	for _, r := range robots {
		router.POST(fmt.Sprintf("/incoming/%s", r.Alias), buildHandler(*r))
	}

	return &Server{
		Engine: router,
	}, nil
}

func incomingHandler(c *gin.Context, r *robot.Robot) {
	body, err := c.GetRawData()
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	params := make(robot.Map)
	if err := json.Unmarshal(body, &params); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	msg, ok := r.MatchMsg(body)
	if !ok {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	resp, err := r.BuildReply(msg.Build(params))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Data(resp.StatusCode(), resp.Header().Get("Content-Type"), resp.Body())
}
