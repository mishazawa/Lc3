include .env
export

build:
	go build -o bin/$(BASE) .

run:
	go run .

exec:
	go run . -exec=data/2048.obj
