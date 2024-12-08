package bot

import (
	"fmt"

	"github.com/go-joe/joe"
	"github.com/samber/lo"
)

type Commands struct {
	bot *joe.Bot
}

func NewCommands(bot *joe.Bot) *Commands {
	return &Commands{bot: bot}
}

func (c *Commands) Handle(command string, chat *Chat) {
	switch command {
	case "get_webhook":
		webhookID := lo.RandomString(16, lo.LowerCaseLettersCharset)
		if err := c.bot.Store.Set(webhookID, chat.Channel); err != nil {
			c.bot.Logger.Error(err.Error())
			return
		}

		webhook := fmt.Sprintf("http://localhost:8080/api/webhooks/%s", webhookID)
		if chat.Channel == chat.ChatID {
			c.bot.Say(chat.ChatID, fmt.Sprintf("Your webhook is: %s", webhook))
			return
		}
		c.bot.Say(chat.ChatID, fmt.Sprintf("Your webhook is: %s", webhook))
		c.bot.Say(chat.Channel, "Congratulations! Your webhook has been sent to you.")
	case "reset_webhook":
		
	case "help":

	}
}
