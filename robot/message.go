package robot

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var tplArgExp = regexp.MustCompile(`{{(\s*\$\S+\s*)}}`)

type Map map[string]interface{}

type Message struct {
	Regexp   string `yaml:"regexp"`
	Template string `yaml:"template"`

	Exp *regexp.Regexp
}

type variable struct {
	full string
	name string
}

func BuildMessage(tpl string, params Map) string {
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

func BuildPostBody(bodyTpl string, message string) *bytes.Buffer {
	return bytes.NewBufferString(strings.Replace(bodyTpl, "$template", message, -1))
}

func extractArgs(params Map, key string) string {
	key = strings.Replace(key, "$", "", -1)
	keys := strings.Split(key, ".")
	for index, k := range keys {
		if index == len(keys)-1 {
			return fmt.Sprintf("%v", params[k])
		}

		if nextParams, ok := params[k].(map[string]interface{}); ok {
			params = nextParams
		}
	}

	return ""
}
