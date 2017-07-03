REPO=malice
NAME=windows-defender
VERSION=$(shell cat VERSION)

all: build size test

dev:
	docker build -t $(REPO)/$(NAME):dev -f Dockerfile.dev .

build:
	docker build -t $(REPO)/$(NAME):$(VERSION) .

size:
	sed -i.bu 's/docker image-.*-blue/docker image-$(shell docker images --format "{{.Size}}" $(REPO)/$(NAME):$(VERSION))-blue/' README.md

tags:
	docker images --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}" $(REPO)/$(NAME)

tar:
	docker save $(REPO)/$(NAME):$(VERSION) -o wdef.tar

test:
	docker run --init --rm $(ORG)/$(NAME):$(VERSION) --help
	test -f befb88b89c2eb401900a68e9f5b78764203f2b48264fcc3f7121bf04a57fd408 || wget https://github.com/maliceio/malice-av/raw/master/samples/befb88b89c2eb401900a68e9f5b78764203f2b48264fcc3f7121bf04a57fd408
	docker run --init --rm -v $(PWD):/malware $(ORG)/$(NAME):$(VERSION) -t befb88b89c2eb401900a68e9f5b78764203f2b48264fcc3f7121bf04a57fd408 > SAMPLE.md
	docker run --init --rm -v $(PWD):/malware $(ORG)/$(NAME):$(VERSION) -V befb88b89c2eb401900a68e9f5b78764203f2b48264fcc3f7121bf04a57fd408 > results.json
	cat results.json | jq .
	rm befb88b89c2eb401900a68e9f5b78764203f2b48264fcc3f7121bf04a57fd408

vagrant:
	@vagrant up
	@vagrant ssh

circle:
	http https://circleci.com/api/v1.1/project/github/${REPO} | jq '.[0].build_num' > .circleci/build_num
	http "$(shell http https://circleci.com/api/v1.1/project/github/${REPO}/$(shell cat .circleci/build_num)/artifacts${CIRCLE_TOKEN} | jq '.[].url')" > .circleci/SIZE
	sed -i.bu 's/docker%20image-.*-blue/docker%20image-$(shell cat .circleci/SIZE)-blue/' README.md

.PHONY: build size tags test tar circle vagrant
