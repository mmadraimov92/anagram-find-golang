FILES = src/*
BINARY = anagram-find

VERSION=`git describe --tags`
BUILD=`date +%FT%T%z`

all: deps build

build:
	CGO_ENABLED=0 GOOS=linux go build -ldflags='-s -w -extldflags "-static"' -o bin/$(BINARY) $(FILES)

	mkdir -p release
	zip release/$(binary_name)-linux.zip bin/$(binary_name)

test:
	go test ./...

bench:
	go test ./... -benchmem -bench=. -run=^$$

deps:
	go get golang.org/x/text/transform
	go get golang.org/x/text/encoding

clean:
	go clean
	rm -rf profile
	rm -rf src.test