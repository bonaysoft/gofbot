package bot

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"github.com/samber/lo"
	"github.com/spf13/viper"
)

type Webhook struct {
	ID string `json:"id"`
}

func NewWebhook() *Webhook {
	randIDHash := md5.Sum([]byte(lo.RandomString(16, lo.LowerCaseLettersCharset)))
	return &Webhook{ID: hex.EncodeToString(randIDHash[:])}
}

func (w *Webhook) String() string {
	addr := viper.GetString("webhook-host")
	if addr == "" {
		addr = fmt.Sprintf("localhost:%d", viper.GetInt("port"))
	}

	return fmt.Sprintf("%s://%s/api/webhooks/%s", viper.GetString("webhook-scheme"), addr, w.ID)
}
