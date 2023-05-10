package robot

import (
	"bytes"
	"log"
	"regexp"
	"strconv"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/bonaysoft/gofbot/pkg/lark"
)

type Map map[string]any

type Message struct {
	Regexp   string `yaml:"regexp"`
	Template string `yaml:"template"`

	Exp *regexp.Regexp
}

func NewMessage(exp string, template string) *Message {
	return &Message{Regexp: exp, Template: template, Exp: regexp.MustCompile(exp)}
}

func (m *Message) Build(params Map) string {
	funcMap := sprig.TxtFuncMap()
	funcMap["larkUserOpenId"] = func(email string) string {
		resp, err := lark.NewClient().GetOpenId(email)
		if err != nil {
			log.Println(err)
			return ""
		}

		return resp.OpenId
	}

	buf := bytes.NewBufferString("")
	t := template.Must(template.New("msg").Funcs(funcMap).Parse(m.Template))
	if err := t.Execute(buf, params); err != nil {
		return ""
	}

	newMsg := buf.String()
	if strconv.CanBackquote(newMsg) {
		return newMsg
	}

	return strconv.Quote(newMsg)
}
