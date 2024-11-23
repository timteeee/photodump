# https://min.io/docs/minio/linux/reference/minio-server/minio-server.html#id5
export MINIO_ROOT_USER := "admin"
export MINIO_ROOT_PASSWORD := "password"

# MinIO dev config
export STORAGE_ENDPOINT := "127.0.0.1:9000"
export STORAGE_ACCESS_KEY := "VPP0fkoCyBZx8YU0QTjH"
export STORAGE_SECRET_KEY := "iFq6k8RLJw5B0faz0cKCXeQk0w9Q8UdtaFzHuw4J"
export STORAGE_BUCKET := "photodump"

# Postgres dev config
export DB_USER := "user"
export DB_PASSWORD := "password"
export DB_HOST := "127.0.0.1"
export DB_PORT := "5432"
export DB_DATABASE := "postgres"

_default:
	just --list

# start containers, run in watch mode, then stop containers
dev: setup
	gow -c \
		-g ./gow-build.sh \
		-e go,html,js,css,mod \
		-i static \
		run main.go -p 8080 --dev

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
