package robot

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildMessage(t *testing.T) {
	params := Map{
		"name": "saltbo",
		"age":  "53",
		"info": Map{
			"city":  "Beijing",
			"email": "yanbo@email.com",
		},
	}
	tpl := `name: {{ .name }}, age: {{ .age }}, city: {{ .info.city }}, openid: {{ larkEmail2openID .info.email }}`
	msg := NewMessage("", tpl).Build(params)
	assert.Contains(t, msg, params["name"])
	assert.Contains(t, msg, params["age"])
	assert.Contains(t, msg, params["info"].(Map)["city"])
	fmt.Println(msg)

	// bodyTpl := `{"msgtype": "markdown", "content": "$template"}`
	// body := buildPostBody(bodyTpl, msg)
	// assert.Contains(t, body.String(), msg)
}
