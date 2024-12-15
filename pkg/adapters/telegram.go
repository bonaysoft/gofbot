package adapters

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"text/template"

	"github.com/go-joe/joe"
	telegram "github.com/robertgzr/joe-telegram-adapter"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/bonaysoft/gofbot/pkg/bot"
)

type Telegram struct {
}

func NewTelegram() bot.Adapter {
	return &Telegram{}
}

func (tg *Telegram) Name() string {
	return "telegram"
}

func (tg *Telegram) Adapter() joe.Module {
	return telegram.Adapter(viper.GetString("TG_TOKEN"))
}

func (tg *Telegram) GetHandlers(jBot *joe.Bot) []any {
	logger := jBot.Logger.Named(tg.Name())
	tga, ok := jBot.Adapter.(*telegram.TelegramAdapter)
	if !ok {
		panic(fmt.Errorf("BUG: Adapter is %T not *telegram.TelegramAdapter", jBot.Adapter))
	}

	u, err := tga.BotAPI.GetMe()
	if err != nil {
		logger.Error("GetMe", zap.Error(err))
		return nil
	}

	cmd := bot.NewCommands(jBot)
	return []any{
		func(ctx context.Context, ev telegram.ReceiveCommandEvent) {
			chatType := bot.ChatTypeP2P
			if ev.Channel() != strconv.Itoa(ev.From.ID) {
				chatType = bot.ChatTypeGroup
			}
			if chatType != bot.ChatTypeP2P && !strings.Contains(ev.Data.Text, u.UserName) {
				logger.Debug("not talking to me", zap.String("text", ev.Data.Text))
				return
			}
			cmd.Handle(&bot.Command{Cmd: ev.Arg0, Args: ev.Args}, &bot.Chat{Provider: tg.Name(), Channel: ev.Channel(), ChatID: strconv.Itoa(ev.From.ID), ChatType: chatType})
		},
		// cause the joe-telegram-adapter can not identity the command with at name syntax
		// Let's handle this situation by ourselves.
		func(ctx context.Context, ev joe.ReceiveMessageEvent) {
			chatType := bot.ChatTypeP2P
			if ev.Channel != ev.AuthorID {
				chatType = bot.ChatTypeGroup
			}

			if chatType != bot.ChatTypeP2P && !strings.Contains(ev.Text, u.UserName) {
				logger.Debug("not talking to me", zap.String("text", ev.Text))
				return
			}

			command, err := cmd.Parse(ev.Text)
			if err != nil {
				logger.Error("cmd.Parse", zap.Error(err))
				return
			}
			cmd.Handle(command, &bot.Chat{Provider: tg.Name(), Channel: ev.Channel, ChatID: ev.AuthorID, ChatType: chatType})
		},
	}
}

func (tg *Telegram) GetFunMap() template.FuncMap {
	return template.FuncMap{}
}
