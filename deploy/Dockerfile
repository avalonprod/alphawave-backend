FROM golang:1.20-alpine3.18 AS builder

WORKDIR /app

COPY . .

RUN go mod download && go get -u ./...
RUN CGO_ENABLED=0 GOOS=linux go build -o ./.bin/app ./cmd/app/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=0 app/.bin/app .
COPY --from=0 app/configs/ ./configs/ 
COPY --from=0 app/templates/ ./templates/


EXPOSE 8000


CMD ["./app"]