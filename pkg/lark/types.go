package lark

type TokenResponse struct {
	Code           int    `json:"code"`
	Msg            string `json:"msg"`
	AppAccessToken string `json:"app_access_token"`
	Expire         int    `json:"expire"`
}

type Response[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

type OpenIdList struct {
	EmailUsers map[string][]OpenIdItem `json:"email_users"`
}

type OpenIdItem struct {
	OpenId string `json:"open_id,omitempty"`
	UserId string `json:"user_id,omitempty"`
}
