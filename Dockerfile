FROM golang:1.20 as builder

ARG GOPROXY
ENV GOPROXY=$GOPROXY
WORKDIR /app

COPY go.* .
RUN go mod download -x

COPY . .
RUN make build

FROM debian:10
RUN apt update && apt install -y openssh-client ca-certificates

WORKDIR /app
COPY --from=builder /app/build/bin/gofbot /app/gofbot
COPY pigeonnest robots
CMD ["/app/gofbot"]