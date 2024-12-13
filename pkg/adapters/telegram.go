package adapters

import (
	"context"
	"strconv"

	"github.com/go-joe/joe"
	telegram "github.com/robertgzr/joe-telegram-adapter"
	"github.com/spf13/viper"

	"github.com/bonaysoft/gofbot/pkg/bot"
)

type Telegram struct {
}

func NewTelegram() *Telegram {
	return &Telegram{}
}

func (tg *Telegram) Name() string {
	return "telegram"
}

func (tg *Telegram) Adapter() joe.Module {
	return telegram.Adapter(viper.GetString("TG_TOKEN"))
}

func (tg *Telegram) GetHandler(jBot *joe.Bot) any {
	cmd := bot.NewCommands(jBot)
	return func(ctx context.Context, ev telegram.ReceiveCommandEvent) {
		chatType := bot.ChatTypeP2P
		if ev.Channel() != strconv.Itoa(ev.From.ID) {
			chatType = bot.ChatTypeGroup
		}
		cmd.Handle(ev.Arg0, &bot.Chat{Provider: tg.Name(), Channel: ev.Channel(), ChatID: strconv.Itoa(ev.From.ID), ChatType: chatType})
	}
}
