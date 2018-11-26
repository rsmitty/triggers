DOCKER_IMAGE_NAME=rsmitty/gcssource
VERSION ?= $(shell bash -c 'read -p "Enter version tag: " pwd; echo $$pwd')

all: linux docker

.PHONY: vendor
vendor:
	dep ensure -v

linux: vendor
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/gcssource .

docker: docker-build docker-push

docker-build:
	docker build -t $(DOCKER_IMAGE_NAME):$(VERSION) .

docker-push:
	docker push $(DOCKER_IMAGE_NAME):$(VERSION)