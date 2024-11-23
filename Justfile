# https://min.io/docs/minio/linux/reference/minio-server/minio-server.html#id5
export MINIO_ROOT_USER := "admin"
export MINIO_ROOT_PASSWORD := "password"

# Development credentials for MinIO
export ACCESS_KEY := "VPP0fkoCyBZx8YU0QTjH"
export SECRET_KEY := "iFq6k8RLJw5B0faz0cKCXeQk0w9Q8UdtaFzHuw4J"
export BUCKET_NAME := "photodump"

# Postgres dev config
export POSTGRES_USER := "user"
export POSTGRES_PASSWORD := "password"
export POSTGRES_DB := "postgres"

_default:
	just --list

# start containers, run in watch mode, then stop containers
dev: setup
	#!/usr/bin/env bash

	DB_URL="postgresql://{{ POSTGRES_USER }}:{{ POSTGRES_PASSWORD }}@127.0.0.1:5432/{{ POSTGRES_DB }}"

	gow -c \
		-g ./gow-build.sh \
		-e go,html,js,css,mod \
		-i static \
		run main.go -p 8080 --storage 127.0.0.1:9000 --bucket {{ BUCKET_NAME }} --db $DB_URL --dev

	just teardown

# build binary
build:
	./pre-build.sh
	go build -o pd cmd/main.go

# start containers
setup:
	docker compose -f docker/docker-compose.yaml up -d

# stop containers
teardown:
	docker compose -f docker/docker-compose.yaml down -v

# clean go cache, remove data dir
clean:
	go clean
	rm -rf data
