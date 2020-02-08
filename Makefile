export BIN_DIR=bin
export PORT=8000
export MACHINE_ID=1
export EXTERNAL_URL=http://0.0.0.0:8080/
export ELASTIC_URL=http://0.0.0.0:9200/
export KAFKA_URL=0.0.0.0:9092
export IMAGE_NAME=freundallein/resender:latest

init:
	git config core.hooksPath .githooks
run:
	go run main.go
test:
	go test -cover ./...
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -a -o $$BIN_DIR/resender
dockerbuild:
	make test
	docker build -t $$IMAGE_NAME -f Dockerfile .
distribute:
	make test
	echo "$$DOCKER_PASSWORD" | docker login -u "$$DOCKER_USERNAME" --password-stdin
	docker build -t $$IMAGE_NAME .
	docker push $$IMAGE_NAME