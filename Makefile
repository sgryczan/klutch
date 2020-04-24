VERSION=0.0.1
IMAGE_NAME=sgryczan/klutch

all: test build-web build-task

test: test-web test-task test-common

test-web:
	cd web && make test

test-task:
	cd task && make test

test-common:
	cd common && make test

build-web:
	docker build -f Dockerfile.web --build-arg VERSION=${VERSION} -t ${IMAGE_NAME}:web-v${VERSION} .

run-server:
	docker run -ti --rm --network="host" server /bin/sh

build-task:
	docker build -f Dockerfile.task --build-arg VERSION=${VERSION} -t ${IMAGE_NAME}:task-v${VERSION} .

push:
	docker push ${IMAGE_NAME}:web-v${VERSION}
	docker push ${IMAGE_NAME}:task-v${VERSION}

.PHONY: test build-web build-task push