package robot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRobot(t *testing.T) {
	_, e := newRobot("../robots/wxwork4gitlab.yaml")
	assert.NoError(t, e)
}
