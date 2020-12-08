IMAGE ?= tweedproject/kernel
NAMESPACE ?= tweed

VERSION ?=
BUILD   ?= $(shell ./build/build-number)
LDFLAGS := -X main.Version="$(VERSION)" -X main.BuildNumber="$(BUILD)"

.PHONY: test

default:
	go fmt . ./api ./cmd/tweed
	go build -ldflags="$(LDFLAGS)" ./cmd/tweed

docker:
	docker build -t $(IMAGE):edge .
	docker push $(IMAGE):edge

push:
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
