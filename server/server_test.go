package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildMessage(t *testing.T) {
	params := Map{
		"name": "saltbo",
		"age":  "53",
		"info": Map{
			"city": "Beijing",
		},
	}
	tpl := `name: {{$name}}, age: {{ $age }}, city: {{ $info.city }}`
	msg := buildMessage(tpl, params)
	assert.Contains(t, msg, params["name"])
	assert.Contains(t, msg, params["age"])
	assert.Contains(t, msg, params["info"].(Map)["city"])

	bodyTpl := `{"msgtype": "markdown", "content": "$template"}`
	body := buildPostBody(bodyTpl, msg)
	assert.Contains(t, body.String(), msg)
}

func performRequest(r http.Handler, method, path string, body io.Reader) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, body)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

type testServer struct {
}

func (ts *testServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b, e := ioutil.ReadAll(r.Body)
	if e != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, err := w.Write(b); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func TestServer_Run(t *testing.T) {
	ts := httptest.NewServer(&testServer{})
	defer ts.Close()

	robots, err := loadRobots("../robots")
	assert.NoError(t, err)

	r, e := New(robots)
	assert.NoError(t, e)
	for _, robot := range robots {
		// reset the hook to the test server URL
		robot.WebHook = ts.URL

		// RUN
		body := bytes.NewBufferString(`{"name": "saltbo", "sex": "man", "info":{"city": "beijing"}}`)
		w := performRequest(r, "POST", fmt.Sprintf("/incoming/%s", robot.Alias), body)

		// TEST
		assert.Equal(t, http.StatusOK, w.Code)
	}
}
