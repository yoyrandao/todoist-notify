FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY . /app

RUN CGO_ENABLED=0 GOOS=linux go generate cmd/telegram-bot/main.go \
  && go build -o dist/main -a -installsuffix cgo cmd/telegram-bot/main.go

FROM alpine

COPY --from=builder /app/dist/main /app/

WORKDIR /app

CMD ["./main"]