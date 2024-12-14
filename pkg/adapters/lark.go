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
	larkevent "github.com/larksuite/oapi-sdk-go/v3/event"
	larkcontact "github.com/larksuite/oapi-sdk-go/v3/service/contact/v3"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	larka "github.com/saltbo/joe-lark-adapter"
	"github.com/spf13/viper"

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

		cmd.Handle(ev.Text, &bot.Chat{Provider: l.Name(), Channel: ev.Channel, ChatID: ev.AuthorID, ChatType: makeChatType(lEv.Message.ChatType)})
	}
}

func (l *Lark) GetFunMap() template.FuncMap {
	return map[string]any{
		"larkEmail2openID":    l.larkEmail2OpenID,
		"larkAtTransform":     l.larkAtTransform,
		"larkAtTransform4All": l.larkAtTransform4All,
	}
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
