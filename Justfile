_default:
	just --list

dev:
	gow -g ./build.sh -e go,html,js,css,mod -c run ./... -p 8080

build:
	bunx tailwindcss -i css/input.css -o public/css/styles.css
	go build -o bin/ cmd/main.go
