GIT_TAG_DESC=$(shell git describe --tags)
TIMESTAMP=$(shell date --rfc-3339=seconds | sed 's/ /T/')
GIT_HASH=$(shell git rev-parse HEAD)

GO_TAGS = -mod vendor

.PHONY: build clean 

.ONESHELL:

build:
	export CGO_ENABLE=0
	go build $(GO_TAGS) -ldflags '-X main.Version=${GIT_TAG_DESC}' ocpp.go 

xz:
	tar cfJ .xz ocpp

clean:
	rm -rf server ocpp.xz
