.PHONY: test

BUILD_TIME := $(shell date "+%F %T")
MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
MKFILE_DIR := $(dir $(MKFILE_PATH))
TARGET = ${MKFILE_DIR}build/gofbot

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
	cd ${MKFILE_DIR} &&  go build -v -ldflags "-s -w    \
		-X main.repo=${GIT_REPO_INFO}					\
		-X main.commit=${COMMIT}						\
		-X main.version=${RELEASE}						\
		-X 'main.buildTime=${BUILD_TIME}'				\
		" -o ${TARGET} ${MKFILE_DIR}server
	@echo "-------------- version detail ---------------"
	@${TARGET} -v

test:
	go test -coverprofile=coverage.txt -covermode=atomic ./...
	go tool cover --func=coverage.txt

covhtml:
	go tool cover -html=coverage.txt

run: build
	${TARGET}

clean:
	rm -rf ${MKFILE_DIR}build
