.PHONY: build run clean all test test-coverage

all: build

build:
	go build -o azenv main.go

run: build
	./azenv

clean:
	rm -f azenv

deps:
	go mod download
