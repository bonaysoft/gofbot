package adapters

import (
	"context"
	"encoding/json"

	"github.com/go-joe/joe"
	larkevent "github.com/larksuite/oapi-sdk-go/v3/event"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	lark "github.com/saltbo/joe-lark-adapter"
	"github.com/spf13/viper"

	"github.com/bonaysoft/gofbot/pkg/bot"
)

type Lark struct {
}

func NewLark() *Lark {
	return &Lark{}
}

func (l *Lark) Name() string {
	return "lark"
}

func (l *Lark) Adapter() joe.Module {
	return lark.Adapter(viper.GetString("LARK_APPID"), viper.GetString("LARK_SECRET"))
}

func (l *Lark) GetHandler(jBot *joe.Bot) any {
	cmd := bot.NewCommands(jBot)
	return func(ctx context.Context, ev joe.ReceiveMessageEvent) {
		larkEvent := &larkevent.EventV2Body{
			Event: &larkim.P2MessageReceiveV1Data{},
		}
		if err := json.Unmarshal(ev.Data.([]byte), larkEvent); err != nil {
			return
		}
		lEv, ok := larkEvent.Event.(*larkim.P2MessageReceiveV1Data)
		if !ok {
			return
		}

		cmd.Handle(ev.Text, &bot.Chat{Channel: ev.Channel, ChatID: ev.AuthorID, ChatType: makeChatType(lEv.Message.ChatType)})
	}
}

var chatTypeMapping = map[string]bot.ChatType{
	"p2p":         bot.ChatTypeP2P,
	"group":       bot.ChatTypeGroup,
	"topic_group": bot.ChatTypeChannel,
}

func makeChatType(chatType *string) bot.ChatType {
	if chatType == nil {
		return bot.ChatTypeP2P
	}

	return chatTypeMapping[*chatType]
}
