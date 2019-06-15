package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"regexp"
	"strings"
)

type Robot struct {
	Name     string     `yaml:"name"`
	Alias    string     `yaml:"uuid"`
	WebHook  string     `yaml:"webhook"`
	BodyTpl  string     `yaml:"bodytpl"`
	Messages []*Message `yaml:"messages"`
}

type Message struct {
	Regexp   string `yaml:"regexp"`
	Template string `yaml:"template"`

	Exp *regexp.Regexp
}

func newRobot(yamlPath string) (*Robot, error) {
	yamlFile, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		return nil, err
	}

	robot := new(Robot)
	if err := yaml.UnmarshalStrict(yamlFile, robot); err != nil {
		return nil, err
	}

	nameHash := md5.Sum([]byte(robot.Name))
	robot.Alias = hex.EncodeToString(nameHash[:])
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
