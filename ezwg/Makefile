.PHONY: all build

all: build

build:
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-s -w" -o dist/ezwg ./cmd/ezwg.go
	upx dist/ezwg
