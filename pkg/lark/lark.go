package lark

import (
	"errors"
	"fmt"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	hc *resty.Client

	cfg     *Config
	apiBase string
}

func NewClient() *Client {
	return &Client{hc: resty.New(), apiBase: "https://open.feishu.cn/open-apis", cfg: NewConfig()}
}

func (c *Client) GetAccessToken() (*TokenResponse, error) {
	var result TokenResponse
	resp, err := c.hc.R().SetBody(c.cfg).SetResult(&result).Post(fmt.Sprintf("%s/auth/v3/app_access_token/internal", c.apiBase))
	if err != nil {
		return nil, fmt.Errorf("error getting app_access_token: %v", err)
	}

	if resp.IsError() {
		return nil, errors.New(resp.String())
	}

	return &result, nil
}

func (c *Client) getToken() string {
	// TODO: get token from cache

	tr, err := c.GetAccessToken()
	if err != nil {
		return ""
	}

	// TODO: set cache
	return tr.AppAccessToken
}

func (c *Client) GetOpenId(email string) (*OpenIdItem, error) {
	var result Response[OpenIdList]
	resp, err := c.hc.R().SetAuthToken(c.getToken()).SetResult(&result).Get(fmt.Sprintf("%s/user/v1/batch_get_id?emails=%s", c.apiBase, email))
	if err != nil {
		return nil, fmt.Errorf("error getting openid: %v", err)
	}

	if resp.IsError() {
		return nil, errors.New(resp.String())
	}

	emailUsers, ok := result.Data.EmailUsers[email]
	if !ok || len(emailUsers) == 0 {
		return nil, fmt.Errorf("not found openid for the email: %s", email)
	}

	return &emailUsers[0], nil
}
