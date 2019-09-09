package robot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildMessage(t *testing.T) {
	params := Map{
		"name": "saltbo",
		"age":  "53",
		"info": map[string]interface{}{
			"city": "Beijing",
		},
	}
	tpl := `name: {{$name}}, age: {{ $age }}, city: {{ $info.city }}`
	msg := BuildMessage(tpl, params)
	assert.Contains(t, msg, params["name"])
	assert.Contains(t, msg, params["age"])
	assert.Contains(t, msg, params["info"].(map[string]interface{})["city"])

	bodyTpl := `{"msgtype": "markdown", "content": "$template"}`
	body := BuildPostBody(bodyTpl, msg)
	assert.Contains(t, body.String(), msg)
}
