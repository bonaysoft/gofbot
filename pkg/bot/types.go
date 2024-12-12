package bot

import "github.com/go-joe/joe"

type ChatType int

const (
	ChatTypeP2P ChatType = iota + 1
	ChatTypeGroup
	ChatTypeChannel
)

type Chat struct {
	Channel  string
	ChatID   string
	ChatType ChatType
}

type Adapter interface {
	Name() string
	Adapter() joe.Module
	GetHandler(bot *joe.Bot) any
}
