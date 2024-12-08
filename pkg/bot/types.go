package bot

import "github.com/go-joe/joe"

type Chat struct {
	Channel string
	ChatID  string
}

type Adapter interface {
	Name() string
	Adapter() joe.Module
	GetHandler(bot *joe.Bot) any
}
