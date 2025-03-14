package adapters

import (
	"fmt"

	"github.com/bonaysoft/gofbot/pkg/bot"
)

type AdapterConstruct func() bot.Adapter

var adapters = map[string]AdapterConstruct{
	"lark":     NewLark,
	"telegram": NewTelegram,
	"terminal": NewTerminal,
}

func GetAdapter(name string) (bot.Adapter, error) {
	adapter, ok := adapters[name]
	if !ok {
		return nil, fmt.Errorf("%s not exist", name)
	}

	return adapter(), nil
}
