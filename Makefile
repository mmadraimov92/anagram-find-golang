FILES = src/*
BINARY = anagram-find

VERSION=`git describe --tags`
BUILD=`date +%FT%T%z`

all: build

build:
	CGO_ENABLED=0 GOOS=darwin go build -ldflags='-s -w -extldflags "-static"' -o bin/$(BINARY) $(FILES)

test:
	go test ./...

bench:
	go test ./... -benchmem -bench=. -run=^$$

clean:
	go clean
	rm -rf profile
	rm -rf src.test