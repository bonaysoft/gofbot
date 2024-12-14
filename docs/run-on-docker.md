# Run with docker

## Via docker-compose

```bash
wget https://raw.githubusercontent.com/bonaysoft/gofbot/master/docker-compose.yml
wget https://raw.githubusercontent.com/bonaysoft/gofbot/master/catalog/gitlab-via-any/note.yaml
docker-compose up -d
```

## Via docker

```bash
docker run -d -p 8080:8080 -e "TG_TOKEN=xxx" bonaysoft/gofbot
```