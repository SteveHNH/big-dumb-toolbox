.PHONY: build clean test

build:
	go build -o bdt .

clean:
	rm -f bdt

test:
	go test ./...

install: build
	install -m 755 bdt $(HOME)/bin/

all: build