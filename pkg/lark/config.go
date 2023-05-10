package lark

import "github.com/spf13/viper"

type Config struct {
	AppID     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
}

func NewConfig() *Config {
	return &Config{
		AppID:     viper.GetString("lark_app_id"),
		AppSecret: viper.GetString("lark_app_secret"),
	}
}
