package lark

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"text/template"
)

func FuncMap() template.FuncMap {
	return map[string]any{
		"larkEmail2openID":    larkEmail2OpenID,
		"larkAtTransform":     larkAtTransform,
		"larkAtTransform4All": larkAtTransform4All,
	}
}

// larkEmail2OpenID 通过邮箱获取openid
func larkEmail2OpenID(email string) string {
	resp, err := NewClient().GetOpenId(email)
	if err != nil {
		log.Println(err)
		return ""
	}

	return resp.OpenId
}

// larkAtTransform 将一个@username文本转换为飞书的at标签
func larkAtTransform(username, emailSuffix string) string {
	openid := larkEmail2OpenID(strings.TrimPrefix(username, "@") + emailSuffix)
	if openid == "" {
		return username
	}

	return fmt.Sprintf("<at id=%s></at>", openid)
}

var atRex = regexp.MustCompile(`@(\w+)`)

// larkAtTransform4All 将一段文本中的所有@username替换为飞书的at标签
func larkAtTransform4All(content, emailSuffix string) string {
	users := atRex.FindAllString(content, -1)
	for _, user := range users {
		content = strings.ReplaceAll(content, user, larkAtTransform(user, emailSuffix))
	}

	return content
}
