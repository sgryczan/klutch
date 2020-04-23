VERSION=0.0.0
IMAGE_NAME=sgryczan/klutch

all: build-web build-task

build-web:
	docker build -f Dockerfile.web -t ${IMAGE_NAME}:web-v${VERSION} .

run-server:
	docker run -ti --rm --network="host" server /bin/sh

build-task:
	docker build -f Dockerfile.task -t ${IMAGE_NAME}:task-v${VERSION} .

push:
	docker push ${IMAGE_NAME}:web-v${VERSION}
	docker push ${IMAGE_NAME}:task-v${VERSION}

.PHONY: build-web build-task push