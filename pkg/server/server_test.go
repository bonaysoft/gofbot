package server

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bonaysoft/gofbot/pkg/robot"
	"github.com/stretchr/testify/assert"
)

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

	robots, err := robot.Load("../robots")
	assert.NoError(t, err)

	s, e := New(robots)
	assert.NoError(t, e)
	for _, r := range robots {
		// reset the hook to the test server URL
		r.WebHook = ts.URL

		// RUN
		body := bytes.NewBufferString(`{"name": "saltbo", "sex": "man", "info":{"city": "beijing"}}`)
		w := performRequest(s, "POST", fmt.Sprintf("/incoming/%s", r.Alias), body)

		// TEST
		assert.Equal(t, http.StatusOK, w.Code)
	}
}
