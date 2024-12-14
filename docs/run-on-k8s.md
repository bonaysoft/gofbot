# Run with Kubernetes

## Installation

```bash
helm add repo bonaysoft https://bonaysoft.github.io/helm-charts
helm install gofbot bonaysoft/gofbot
```

## Usage

```bash
kubectl apply -f https://raw.githubusercontent.com/bonaysoft/gofbot/master/catalog/gitlab-via-any/general.yaml
kubectl apply -f https://raw.githubusercontent.com/bonaysoft/gofbot/master/catalog/gitlab-via-lark/note.yaml
```