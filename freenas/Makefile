REGISTRY = registry.wheat-dn42.net
IMAGE_PROVISONER = $(REGISTRY)/freenas-provisioner

VERSION := latest

all build:
	#GOOS=linux go install -v ./freenas-provisioner
	CGO_ENABLED=0 GOOS=linux go build -o freenas-provisioner .
.PHONY: all build

container: build quick-container
.PHONY: container

quick-container:
	cp freenas-provisioner deploy/docker/freenas-provisioner
	docker build -t $(IMAGE_PROVISONER):$(VERSION) -f deploy/docker/Dockerfile deploy/docker
.PHONY: quick-container

push: container
	docker push $(IMAGE_PROVISONER):$(VERSION)
.PHONY: push

clean:
	rm -f freenas-provisioner
	rm -f deploy/docker/freenas-provisioner
.PHONY: clean
