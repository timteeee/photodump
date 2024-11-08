_default:
	just --list

dev:
	gow -e go,html,css,mod run cmd/main.go -p 8080

build:
	go build -o bin/ cmd/main.go
