IMAGE ?= tweedproject/kernel
NAMESPACE ?= tweed

VERSION ?=
BUILD   ?= $(shell ./build/build-number)
LDFLAGS := -X main.Version="$(VERSION)" -X main.BuildNumber="$(BUILD)"
PWD = $(shell pwd)

# if VERSION isn't specified (i.e. ""), then
# use the tag 'latest' per Docker pratices.
DOCKER_TAG ?= $(VERSION)
ifeq ($(DOCKER_TAG),)
  DOCKER_TAG := latest
endif

.PHONY: test default docker deploy.yml retire deploy redeploy push test bg-tests

default:
	go mod vendor
	go fmt . ./api ./cmd/tweed
	go build -mod=vendor -ldflags="$(LDFLAGS)" ./cmd/tweed

deploy.yml:
	@IMAGE=$(IMAGE) VERSION=$(DOCKER_TAG) NAMESPACE=$(NAMESPACE) envsubst <env/dev/k8s.yml

retire:
	kubectl delete ns $(NAMESPACE) || true

deploy:
	IMAGE=$(IMAGE) VERSION=$(DOCKER_TAG) NAMESPACE=$(NAMESPACE) envsubst <env/dev/k8s.yml | kubectl apply  -f -
	@echo
	@echo "    Now run \`source env/dev/envrc\` to update your current shell environment."
	@echo

redeploy: retire deploy

docker: default
	go mod vendor
	docker build -t $(IMAGE):edge .

push: default
	@echo "Checking that VERSION was defined in the calling environment"
	@test -n "$(VERSION)"
	@echo "OK.  VERSION=$(VERSION)"
	docker build -t $(IMAGE):$(VERSION) .
	docker push $(IMAGE):$(VERSION)
	docker tag $(IMAGE):$(VERSION) $(IMAGE):latest
	for V in $(VERSION) $(shell echo "$(VERSION)" | sed -e 's/\.[^.]*$$//') $(shell echo "$(VERSION)" | sed -e 's/\..*$$//'); do \
		docker tag $(IMAGE):$(VERSION) $(IMAGE):$$V; \
		docker push $(IMAGE):$$V; \
	done

test:
	./test/the all

bg-tests:
	docker build -t tweed-unit -f Dockerfile.test .
	docker run --rm -it --privileged \
	           --mount type=bind,source=$(PWD),target=/tweed,consistency=cached \
	           tweed-unit:latest \
	           ginkgo watch ./...
