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

- If you want One robot connect to multiple chat platform at the same time, you should start up multiple robots.
- If you cannot tolerate message loss, you should use it with the event bus.

## Getting Started

### Step1: Install and boot your bot

> If you don't have a bot, you should create one
> first, [see here](https://core.telegram.org/bots#3-how-do-i-create-a-bot).

```bash
wget https://raw.githubusercontent.com/bonaysoft/gofbot/master/docker-compose.yml
echo "TG_TOKEN=xxxx" > .env
docker-compose up
```

### Step2: Write your Message for your bot

```yaml
apiVersion: github.com/bonaysoft/gofbot/v1alpha1
kind: Message
metadata:
  name: gl-push
spec:
  selector:
    matchLabels:
      type: xxCreate

  reply:
    text: |-
      This is a test message.
      name: {{ .name }}
      intro: {{ .intro }}
      hello: {{ .other.key }}
```

### Step2: Get the webhook

Chat with your robot and get the webhook address.

Just send the command `/get_webhook` then the robot will send you the webhook address.

### Step4: Post the webhook

Post any body to the webhook address, like this:

```bash
curl -X POST -H "Content-Type: application/json" -d '{"type": "xxCreate", "name": "xxx", "intro": "this is a test", "other": {"key": "world"}}' $WEBHOOK_URL
```

## Supports

- [x] Telegram
- [ ] Slack
- [ ] RocketChat
- [ ] Mattermost
- [ ] Discord
- [x] Lark
- [ ] DingTalk
- [ ] WeCom

## Need Helps (Welcome PRs)

- Implement more adapters
- Build more Messages into the [catalog](catalog)
- Rich the documents

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available,
see the [tags on this repository][tags].

## Authors

- **Ambor** - *Initial work* - [saltbo](https://github.com/saltbo)

See also the list of [contributors][contributors] who participated in this project.
