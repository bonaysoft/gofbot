package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
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

func loadRobots(robotsPath string) ([]*Robot, error) {
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
