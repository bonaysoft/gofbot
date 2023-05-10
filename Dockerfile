FROM golang:1.20 as builder

ARG GOPROXY
ENV GOPROXY=$GOPROXY
WORKDIR /app

COPY go.* .
RUN go mod download -x

COPY . .
RUN make build

FROM debian:10

WORKDIR /app

COPY --from=builder /app/build/bin/gofbot /app/gofbot
COPY robots robots
CMD ["/app/gofbot"]