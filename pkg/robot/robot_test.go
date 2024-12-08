package robot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRobot(t *testing.T) {
	_, e := NewRobot("./testdata/example.yaml")
	assert.NoError(t, e)
}
