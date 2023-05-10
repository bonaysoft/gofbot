# gofbot

[![build](https://github.com/bonaysoft/gofbot/actions/workflows/build.yml/badge.svg)](https://github.com/bonaysoft/gofbot/actions/workflows/build.yml)
[![codecov](https://codecov.io/gh/bonaysoft/gofbot/branch/master/graph/badge.svg)](https://codecov.io/gh/bonaysoft/gofbot)
[![Go Report Card](https://goreportcard.com/badge/github.com/bonaysoft/gofbot)](https://goreportcard.com/report/github.com/bonaysoft/gofbot)

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