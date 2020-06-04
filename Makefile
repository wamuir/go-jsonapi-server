.PHONY: build create rm run start stop

NS = github.com/wamuir
IMAGE_NAME = go-jsonapi-server
VERSION = latest

CONTAINER_NAME = go-jsonapi-server
CONTAINER_INSTANCE = default

ENV  = 
NETWORK = --network cluster
PORTS = --publish 8080:8080
VOLUMES =

build:
	docker build -t $(NS)/$(IMAGE_NAME):$(VERSION) -f Dockerfile .

create:
	docker create --restart unless-stopped --name $(CONTAINER_NAME)-$(CONTAINER_INSTANCE) $(ENV) $(NETWORK) $(PORTS) $(VOLUMES) $(NS)/$(IMAGE_NAME):$(VERSION)

rm:
	docker rm $(CONTAINER_NAME)-$(CONTAINER_INSTANCE)

run:
	docker run -d --restart unless-stopped --name $(CONTAINER_NAME)-$(CONTAINER_INSTANCE) $(ENV) $(NETWORK) $(PORTS) $(VOLUMES) $(NS)/$(IMAGE_NAME):$(VERSION)

start:
	docker start $(CONTAINER_NAME)-$(CONTAINER_INSTANCE)

stop:
	docker stop $(CONTAINER_NAME)-$(CONTAINER_INSTANCE)

reload:
	docker kill -s HUP $(CONTAINER_NAME)-$(CONTAINER_INSTANCE)

default: build
