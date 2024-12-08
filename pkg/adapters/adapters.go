package adapters

import (
	"fmt"

	"github.com/bonaysoft/gofbot/pkg/bot"
)

var adapters = map[string]bot.Adapter{
	"telegram": NewTelegram(),
}

func GetAdapter(name string) (bot.Adapter, error) {
	adapter, ok := adapters[name]
	if !ok {
		return nil, fmt.Errorf("%s not exist", name)
	}

	return adapter, nil
}
