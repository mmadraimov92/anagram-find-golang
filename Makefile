files = main.go anagram.go charset_table.go
binary_name = anagram-find

all: deps build

build:
	CGO_ENABLED=0 GOOS=linux go build -ldflags='-s -w -extldflags "-static"' -o bin/$(binary_name) $(files)
	zip $(binary_name)-linux.zip bin/$(binary_name)

deps:
	go get golang.org/x/text/transform
	go get golang.org/x/text/encoding

clean:
	go clean
	rm -rf profile