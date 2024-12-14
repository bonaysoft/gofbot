# gofbot

[![build](https://github.com/bonaysoft/gofbot/actions/workflows/build.yml/badge.svg)](https://github.com/bonaysoft/gofbot/actions/workflows/build.yml)
[![codecov](https://codecov.io/gh/bonaysoft/gofbot/branch/master/graph/badge.svg)](https://codecov.io/gh/bonaysoft/gofbot)
[![Go Report Card](https://goreportcard.com/badge/github.com/bonaysoft/gofbot)](https://goreportcard.com/report/github.com/bonaysoft/gofbot)

A generic forwarding robots for any webhooks.

It was a generic webhook endpoint to call some notifications.

## Features

- Chat with a robot to manage your webhook
- Match and reply templates through declarative file descriptions
- Rich template functions help you implement your own reply templates
- Support multiple chat platforms
- Support multiple webhook platforms

## Non-Features

- One robot connect to multiple chat platform, your should use multiple robots
- Message persistence is not supported and messages are not guaranteed to be lost, If you cannot tolerate message loss,
  you should use it with the event bus

## Getting Started

### Step1: Write your Message

```yaml
apiVersion: github.com/bonaysoft/gofbot/v1alpha1
kind: Message
metadata:
  name: gl-push
spec:
  selector:
    matchLabels:
      chatProvider: lark
      event_type: issue

  reply:
    text: |-
      test push message.
      name: {{ .name }}
      intro: {{ .intro }}
      hello: {{ .other.key }}
```

### Step2: Run the robot

> If you don't have a bot, you should create one
> first, [see here](https://core.telegram.org/bots#3-how-do-i-create-a-bot).

```bash
TG_TOKEN=xxx gofbot run --adapter telegram
```

### Step3: Get the webhook

Chat with your robot and get the webhook address.

Just send the command `/get_webhook` then the robot will send you the webhook address.

### Step4: Post the webhook

Post any body to the webhook address, like this:

```bash
curl -X POST -H "Content-Type: application/json" -d '{"type": "xxCreate", "name": "xxx", "intro": "this is a test", "other": {"key": "world"}}' $WEBHOOK_URL
```

## Supports

-[x] Telegram
-[ ] Slack
-[ ] RocketChat
-[ ] Mattermost
-[ ] Discord
-[x] Lark
-[ ] DingTalk
-[ ] WeCom

