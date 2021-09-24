all: build
	./smash

build:
	go build

test:
	go test

rebuild: test all
	