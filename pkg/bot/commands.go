package bot

import (
	"fmt"
	"strings"

	"github.com/go-joe/joe"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type Command struct {
	Cmd  string
	Args []string
}

type Commands struct {
	bot *joe.Bot
}

func NewCommands(bot *joe.Bot) *Commands {
	return &Commands{bot: bot}
}

func (c *Commands) Parse(text string) (*Command, error) {
	items := lo.Filter(strings.Split(text, " "), func(item string, index int) bool { return !strings.HasPrefix(item, "@") })
	return &Command{Cmd: strings.TrimPrefix(items[0], "/"), Args: items[1:]}, nil
}

func (c *Commands) Handle(command *Command, chat *Chat) {
	switch command.Cmd {
	case "get_webhook":
		c.getWebhook(chat)
	case "reset_webhook":
		c.resetWebhook(chat)
	case "help":
		c.sendHelp(chat)
	}
}

func (c *Commands) sendHelp(chat *Chat) {
	helpText := "1. get_webhook: get one webhook for the current chat.\\n2. reset_webhook: reset the webhook for the current chat."
	c.bot.Say(chat.Channel, helpText)
}

func (c *Commands) getWebhook(chat *Chat) {
	webhook, err := c.getOrCreateWebhook(chat)
	if err != nil {
		c.bot.Logger.Error("getOrCreateWebhook", zap.Error(err))
		return
	}

	if chat.ChatType == ChatTypeP2P {
		c.bot.Say(chat.ChatID, fmt.Sprintf("Your webhook is: %s", webhook))
		return
	}
	c.bot.Say(chat.ChatID, fmt.Sprintf("Your webhook is: %s", webhook))
	c.bot.Say(chat.Channel, "Congratulations! Your webhook has been sent to you.")
}

func (c *Commands) resetWebhook(chat *Chat) {
	webhook, err := c.createWebhook(chat)
	if err != nil {
		c.bot.Logger.Error("createWebhook", zap.Error(err))
		return
	}
	if chat.ChatType == ChatTypeP2P {
		c.bot.Say(chat.ChatID, fmt.Sprintf("Your new webhook is: %s", webhook))
		return
	}

	c.bot.Say(chat.ChatID, fmt.Sprintf("Your new webhook is: %s", webhook))
	c.bot.Say(chat.Channel, "Congratulations! Your webhook has been sent to you.")
}

func (c *Commands) getOrCreateWebhook(chat *Chat) (*Webhook, error) {
	var webhook Webhook
	ok, err := c.bot.Store.Get(fmt.Sprintf("CH_%s", chat.Channel), &webhook)
	if err != nil {
		return nil, err
	} else if ok {
		return &webhook, nil
	}

	return c.createWebhook(chat)
}

func (c *Commands) createWebhook(chat *Chat) (*Webhook, error) {
	webhook := NewWebhook()
	if err := c.bot.Store.Set(fmt.Sprintf("WH_%s", webhook.ID), chat); err != nil {
		return nil, err
	}

	if err := c.bot.Store.Set(fmt.Sprintf("CH_%s", chat.Channel), webhook); err != nil {
		return nil, err
	}

	return webhook, nil
}
