BINARY_NAME=news-aggregator
.DEFAULT_GOAL := run


build:
GOARCH=amd64 GOOS=linux go build -o ./target/${BINARY_NAME}-linux main.go
