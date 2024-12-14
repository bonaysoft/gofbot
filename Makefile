.PHONY: test build

BUILD_TIME := $(shell date "+%F %T")
MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
MKFILE_DIR := $(dir $(MKFILE_PATH))
TARGET = ${MKFILE_DIR}build/bin/gofbot

RELEASE?=$(shell git describe --tags)
GIT_REPO_INFO=$(shell git config --get remote.origin.url)

ifndef COMMIT
  COMMIT := git-$(shell git rev-parse --short HEAD)
endif

default: build

mod:
	@echo "-------------- download the modules ---------------"
	@go mod download

build:
	@echo "-------------- building the program ---------------"
	cd ${MKFILE_DIR} && CGO_ENABLED=0 go build -v -ldflags "-s -w    \
		-X github.com/bonaysoft/gofbot/cmd.repo=${GIT_REPO_INFO}					\
		-X github.com/bonaysoft/gofbot/cmd.commit=${COMMIT}						\
		-X github.com/bonaysoft/gofbot/cmd.release=${RELEASE}						\
		" -o ${TARGET}
	@echo "-------------- version detail ---------------"
	@${TARGET} -v

test:
	go test -coverprofile=.coverprofile -covermode=atomic ./...
	go tool cover --func=.coverprofile

covhtml:
	go tool cover -html=.coverprofile

run: build
	${TARGET}

clean:
	rm -rf ${MKFILE_DIR}build
