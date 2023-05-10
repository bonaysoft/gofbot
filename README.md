# gofbot

[![build](https://github.com/bonaysoft/gofbot/actions/workflows/build.yml/badge.svg)](https://github.com/bonaysoft/gofbot/actions/workflows/build.yml)
[![codecov](https://codecov.io/gh/bonaysoft/gofbot/branch/master/graph/badge.svg)](https://codecov.io/gh/bonaysoft/gofbot)
[![Go Report Card](https://goreportcard.com/badge/github.com/bonaysoft/gofbot)](https://goreportcard.com/report/github.com/bonaysoft/gofbot)

A generic forwarding robots based on the URL called.

It was a generic webhook endpoint to call some notifications.

## 内置函数

- larkEmail2openID: 通过邮箱获取openid
- larkAtTransform: 将一个@username文本转换为飞书的at标签
- larkAtTransform4All: 将一段文本中的所有@username替换为飞书的at标签

### 用法

```gotemplate

<at id={{ larkEmail2openID (printf "%s@%s" .user.username "company.com") }}></at>
```

## 环境变量

***使用飞书相关函数时必须设置以下环境变量***

- LARK_APP_ID=cli_a4d34d4xxxxxxxx
- LARK_APP_SECRET=y9uUaXQqz71vf0cmxyxxxxxxxxxxxx

