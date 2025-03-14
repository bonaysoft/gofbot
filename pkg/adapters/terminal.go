package adapters

import (
	"context"
	"text/template"

	"github.com/go-joe/joe"

	"github.com/bonaysoft/gofbot/pkg/bot"
)

type Terminal struct {
}

func NewTerminal() bot.Adapter {
	return &Terminal{}
}

func (c *Terminal) Name() string {
	return "terminal"
}

func (c *Terminal) Adapter() joe.Module {
	return joe.ModuleFunc(func(config *joe.Config) error {
		config.SetAdapter(joe.NewCLIAdapter("terminal", config.Logger("terminal")))
		return nil
	})
}

func (c *Terminal) GetHandlers(jBot *joe.Bot) []any {
	cmd := bot.NewCommands(jBot)
	return []any{
		func(ctx context.Context, ev joe.ReceiveMessageEvent) {
			cmd.Handle(&bot.Command{Cmd: ev.Text}, &bot.Chat{Provider: c.Name(), ChatID: ev.ID})
		},
	}
}

func (c *Terminal) GetFunMap() template.FuncMap {
	return template.FuncMap{}
}
