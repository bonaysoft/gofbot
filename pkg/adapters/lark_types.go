package adapters

import (
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/samber/lo"
)

type LarkBotInfoResponse struct {
	larkcore.CodeError
	Bot LarkBotInfo `json:"bot"`
}

type LarkBotInfo struct {
	ActivateStatus int           `json:"activate_status"`
	AppName        string        `json:"app_name"`
	AvatarUrl      string        `json:"avatar_url"`
	IpWhiteList    []interface{} `json:"ip_white_list"`
	OpenId         string        `json:"open_id"`
}

// isTalk2me checks if the bot is mentioned in the given list of mentions.
func isTalk2me(mentions []*larkim.MentionEvent, botInfo *LarkBotInfo) bool {
	_, ok := lo.Find(mentions, func(item *larkim.MentionEvent) bool { return *item.Id.OpenId == botInfo.OpenId })
	return ok
}
