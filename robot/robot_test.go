package robot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRobot(t *testing.T) {
	_, e := newRobot("../deployments/robots/wxwork4gitlab.yaml")
	assert.NoError(t, e)
}
