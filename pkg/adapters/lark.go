package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"text/template"

	"github.com/go-joe/joe"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkevent "github.com/larksuite/oapi-sdk-go/v3/event"
	larkcontact "github.com/larksuite/oapi-sdk-go/v3/service/contact/v3"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	larka "github.com/saltbo/joe-lark-adapter"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/bonaysoft/gofbot/pkg/bot"
	"github.com/bonaysoft/gofbot/pkg/errors"
)

type Lark struct {
	appId, appSecret string

	client *lark.Client
}

func NewLark() bot.Adapter {
	appId, appSecret := viper.GetString("LARK_APPID"), viper.GetString("LARK_SECRET")
	return &Lark{
		appId:     appId,
		appSecret: appSecret,
		client:    lark.NewClient(appId, appSecret),
	}
}

func (l *Lark) Name() string {
	return "lark"
}

func (l *Lark) Adapter() joe.Module {
	return larka.Adapter(l.appId, l.appSecret)
}

func (l *Lark) GetHandlers(jBot *joe.Bot) []any {
	logger := jBot.Logger.Named(l.Name())
	botInfo, err := l.GetMe(context.Background())
	if err != nil {
		panic(err)
	}

	cmd := bot.NewCommands(jBot)
	return []any{
		func(ctx context.Context, ev joe.ReceiveMessageEvent) {
			larkEvent := &larkevent.EventV2Body{Event: &larkim.P2MessageReceiveV1Data{}}
			if err := json.Unmarshal(ev.Data.([]byte), larkEvent); err != nil {
				logger.Error("json.Unmarshal", zap.Error(err))
				return
			}

			lEv, ok := larkEvent.Event.(*larkim.P2MessageReceiveV1Data)
			if !ok {
				return
			}

			chatType := makeChatType(lEv.Message.ChatType)
			if chatType != bot.ChatTypeP2P && !isTalk2me(lEv.Message.Mentions, botInfo) {
				logger.Debug("not talking to me", zap.String("text", ev.Text))
				return
			}
			command, err := cmd.Parse(ev.Text)
			if err != nil {
				logger.Error("cmd.Parse", zap.Error(err))
				return
			}

			cmd.Handle(command, &bot.Chat{Provider: l.Name(), Channel: ev.Channel, ChatID: ev.AuthorID, ChatType: chatType})
		},
	}
}

func (l *Lark) GetFunMap() template.FuncMap {
	return map[string]any{
		"larkEmail2openID":    l.larkEmail2OpenID,
		"larkAtTransform":     l.larkAtTransform,
		"larkAtTransform4All": l.larkAtTransform4All,
	}
}

func (l *Lark) GetMe(ctx context.Context) (*LarkBotInfo, error) {
	resp, err := l.client.Get(ctx, "/open-apis/bot/v3/info", nil, larkcore.AccessTokenTypeTenant)
	if err != nil {
		return nil, err
	}

	var botInfoResp LarkBotInfoResponse
	if err := resp.JSONUnmarshalBody(&botInfoResp, &larkcore.Config{Serializable: &larkcore.DefaultSerialization{}}); err != nil {
		return nil, err
	}

	return &botInfoResp.Bot, nil
}

// larkEmail2OpenID 通过邮箱获取openid
func (l *Lark) larkEmail2OpenID(email string) string {
	ctx := context.Background()
	body := larkcontact.NewBatchGetIdUserReqBodyBuilder().Emails([]string{email}).Build()
	req := larkcontact.NewBatchGetIdUserReqBuilder().UserIdType("open_id").Body(body).Build()
	resp, err := l.client.Contact.User.BatchGetId(ctx, req)
	if err != nil {
		slog.Error("larkEmail2OpenID", errors.With(err))
		return email
	} else if !resp.Success() {
		slog.Error("larkEmail2OpenID", errors.With(resp.CodeError))
		return email
	} else if len(resp.Data.UserList) == 0 {
		slog.Error("larkEmail2OpenID", "error", fmt.Errorf("not found"))
		return email
	}

	if uid := resp.Data.UserList[0].UserId; uid != nil {
		return *uid
	}

	return email
}

// larkAtTransform 将一个@username文本转换为飞书的at标签
func (l *Lark) larkAtTransform(username, emailSuffix string) string {
	openid := l.larkEmail2OpenID(strings.TrimPrefix(username, "@") + emailSuffix)
	if openid == "" {
		return username
	}

	return fmt.Sprintf("<at id=%s></at>", openid)
}

var atRex = regexp.MustCompile(`@(\w+)`)

// larkAtTransform4All 将一段文本中的所有@username替换为飞书的at标签
func (l *Lark) larkAtTransform4All(content, emailSuffix string) string {
	users := atRex.FindAllString(content, -1)
	for _, user := range users {
		content = strings.ReplaceAll(content, user, l.larkAtTransform(user, emailSuffix))
	}

	return content
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
