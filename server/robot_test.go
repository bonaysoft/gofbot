package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRobot(t *testing.T) {
	_, e := newRobot("../robots/wxwork4gitlab.yaml")
	assert.NoError(t, e)
}
