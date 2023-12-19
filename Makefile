.PHONY:
.SILENT:
.DEFAULT_GOAL := run

build:
	go mod download && CGO_ENABLED=0 GOOS=linux go build -o ./.bin/app ./cmd/app/main.go

run: build
	docker compose up --remove-orphans app

stop: 
	docker compose stop && docker compose rm -f

lint: 
	golangci-lint run
