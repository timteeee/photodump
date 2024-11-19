# https://min.io/docs/minio/linux/reference/minio-server/minio-server.html#id5
export MINIO_ROOT_USER := "admin"
export MINIO_ROOT_PASSWORD := "password"

# Development credentials for MinIO
export ACCESS_KEY := "VPP0fkoCyBZx8YU0QTjH"
export SECRET_KEY := "iFq6k8RLJw5B0faz0cKCXeQk0w9Q8UdtaFzHuw4J"
export BUCKET_NAME := "photodump"

_default:
	just --list

# start containers, run in watch mode, then stop containers
dev: setup
	gow -c \
		-g ./gow-build.sh \
		-e go,html,js,css,mod \
		-i internal/assets/public \
		run cmd/main.go -p 8080 --storage 127.0.0.1:9000 --bucket photodump --dev
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
