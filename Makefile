GIT_BRANCH=$(shell git rev-parse --abbrev-ref HEAD 2>/dev/null)
GIT_COMMIT=$(shell git rev-parse --short HEAD)
LDFLAGS=-ldflags "-s -w -X main.GitBranch=${GIT_BRANCH} -X main.GitCommit=${GIT_COMMIT} -X main.BuildDate=`date -u +%Y-%m-%d.%H:%M:%S`"

proto:
	@echo "Make Pinba proto"
	@protoc --gofast_out=. proto/pinba/*.proto

build:
	@[ -d .build ] || mkdir -p .build
	CGO_ENABLED=0 go build ${LDFLAGS} -o .build/proton-server cmd/proton-server/main.go
	@file  .build/proton-server
	@du -h .build/proton-server

deb: build
	@nfpm pkg --target .build/proton-server.deb
	@dpkg-deb -I .build/proton-server.deb

.PHONY: proto