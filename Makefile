REPO=malice
NAME=windows-defender
VERSION=$(shell cat VERSION)

all: build size test

dev:
		docker build -t $(REPO)/$(NAME):$(VERSION) -f Dockerfile.dev .

build:
	docker build -t $(REPO)/$(NAME):$(VERSION) .

size:
	sed -i.bu 's/docker image-.*-blue/docker image-$(shell docker images --format "{{.Size}}" $(REPO)/$(NAME):$(VERSION))-blue/' README.md

tags:
	docker images --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}" $(REPO)/$(NAME)

tar:
	docker save $(REPO)/$(NAME):$(VERSION) -o wdef.tar

test:
	docker run --init --rm $(REPO)/$(NAME):$(VERSION)
	docker run --init --rm $(REPO)/$(NAME):$(VERSION) -V EICAR > results.json
	cat results.json | jq .

.PHONY: build size tags test tar
