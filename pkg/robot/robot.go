package robot

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/samber/lo"
	"gopkg.in/yaml.v3"
)

var hc = resty.New()

type Robot struct {
	Name     string     `yaml:"name"`
	Alias    string     `yaml:"uuid"`
	WebHook  string     `yaml:"webhook"`
	BodyTpl  string     `yaml:"bodytpl"`
	Messages []*Message `yaml:"messages"`
}

func newRobot(yamlPath string) (*Robot, error) {
	yamlFile, err := os.ReadFile(yamlPath)
	if err != nil {
		return nil, err
	}

	robot := new(Robot)
	if err := yaml.Unmarshal(yamlFile, robot); err != nil {
		return nil, err
	}

	nameHash := md5.Sum([]byte(robot.Name))
	robot.Alias = hex.EncodeToString(nameHash[:])
	fmt.Printf("%s => %s\n", robot.Name, robot.Alias)
	errors := make([]string, 0)
	for _, msg := range robot.Messages {
		exp, err2 := regexp.Compile(msg.Regexp)
		if err != nil {
			errors = append(errors, err2.Error())
			continue
		}

		msg.Exp = exp
	}

	if len(errors) != 0 {
		return nil, fmt.Errorf(strings.Join(errors, "\n"))
	}

	return robot, nil
}

func (r *Robot) MatchMsg(body []byte) (*Message, bool) {
	return lo.Find(r.Messages, func(item *Message) bool {
		return item.Exp.Match(body)
	})
}

func (r *Robot) BuildReply(msg string) (*resty.Response, error) {
	body := bytes.NewBufferString(strings.Replace(r.BodyTpl, "$template", msg, -1))
	resp, err := hc.R().SetBody(body).Post(r.WebHook)
	if err != nil {
		return nil, fmt.Errorf("dispatch to robot %s: %v", r.Name, err)
	}

	return resp, nil
}

func findRobots(root string, creator func(filepath string) error) error {
	return filepath.Walk(root, func(filepath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		} else if info.IsDir() {
			return nil
		} else if path.Ext(filepath) != ".yaml" && path.Ext(filepath) != ".yml" {
			return nil
		}

		return creator(filepath)
	})
}

func Load(robotsPath string) ([]*Robot, error) {
	robots := make([]*Robot, 0)
	robotCreator := func(filepath string) error {
		robot, err := newRobot(filepath)
		if err != nil {
			return err
		}

		robots = append(robots, robot)
		return nil
	}

	if err := findRobots(robotsPath, robotCreator); err != nil {
		return nil, err
	}

	if len(robots) == 0 {
		return nil, fmt.Errorf("not found any robot.")
	}

	return robots, nil
}
