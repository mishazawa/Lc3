include .env
export

build:
	go build -o bin/$(BASE) .

run:
	go run .
