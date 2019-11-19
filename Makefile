IMAGE := tweedproject/kernel

VERSION ?= 
BUILD   ?= $(shell ./build/build-number)
LDFLAGS := -X main.Version="$(VERSION)" -X main.BuildNumber="$(BUILD)"

default:
	go fmt . ./api ./cmd/tweed
	go build -ldflags="$(LDFLAGS)" ./cmd/tweed

docker:
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
	./test/the shared
	./test/the dedicated
