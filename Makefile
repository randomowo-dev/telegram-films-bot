.PHONY: run build test prepare

run: build
	docker-compose -f deployments/docker-compose.yml up

build: test
	docker-compose -f deployments/docker-compose.yml build

test: prepare
	go test ./...

prepare:
	go mod tidy
	go mod verify