package adapters

import (
	"context"
	"strconv"

	"github.com/go-joe/joe"
	telegram "github.com/robertgzr/joe-telegram-adapter"

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
	return telegram.Adapter("7887174654:AAGbt87xaM1FTpTWqeq3Z1lo8VMqtN91GpA")
}

func (tg *Telegram) GetHandler(jBot *joe.Bot) any {
	cmd := bot.NewCommands(jBot)
	return func(ctx context.Context, ev telegram.ReceiveCommandEvent) {
		cmd.Handle(ev.Arg0, &bot.Chat{Channel: ev.Channel(), ChatID: strconv.Itoa(ev.From.ID)})
	}
}
