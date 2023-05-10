package lark

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_GetOpenId(t *testing.T) {
	c := NewClient()
	v, err := c.GetOpenId("yanbo@lixiang.com")
	assert.NoError(t, err)
	fmt.Println(v)
}
