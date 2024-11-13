_default:
	just --list

dev:
	gow -c \
		-g ./gow-build.sh \
		-e go,html,js,css,mod \
		-i assets/public \
		run cmd/main.go -p 8080

build:
	./pre-build.sh
	go build -o pd cmd/main.go
