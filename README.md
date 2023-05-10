# gofbot

[![build](https://github.com/bonaysoft/gofbot/actions/workflows/build.yml/badge.svg)](https://github.com/bonaysoft/gofbot/actions/workflows/build.yml)
[![codecov](https://codecov.io/gh/saltbo/gofbot/branch/master/graph/badge.svg)](https://codecov.io/gh/saltbo/gofbot)
[![codebeat badge](https://codebeat.co/badges/e97d3305-de49-4a9c-9ead-1aca942b9e16)](https://codebeat.co/projects/github-com-saltbo-gofbot-master)
[![Go Report Card](https://goreportcard.com/badge/github.com/saltbo/gofbot)](https://goreportcard.com/report/github.com/saltbo/gofbot)

A generic forwarding robots based on the URL called.

It was a generic webhook endpoint to call some notifications.

## Usage

## 内置函数

### larkUserOpenId

#### Required Environment Variables

- LARK_APP_ID=cli_a4d34d4xxxxxxxx
- LARK_APP_SECRET=y9uUaXQqz71vf0cmxyxxxxxxxxxxxx

#### 用法
```gotemplate
<at id={{ larkUserOpenId (printf "%s@%s" .user.username "company.com") }}></at>
```