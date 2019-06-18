FROM golang:latest AS build-env

ENV GOPROXY https://goproxy.io
ENV APP_HOME /app
WORKDIR $APP_HOME

ADD go.mod go.sum Makefile .git $APP_HOME/
RUN make mod

ADD . $APP_HOME
RUN make


# Runing Environment
FROM debian:9

ENV APP_HOME /app
WORKDIR $APP_HOME

ADD robots $APP_HOME/robots
COPY --from=build-env /app/build/gofbot $APP_HOME/gofbot

ENTRYPOINT ["./gofbot"]
