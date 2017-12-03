.PHONY: dev build image test deps clean

CGO_ENABLED=0
COMMIT=`git rev-parse --short HEAD`
APP=cadmus
PACKAGE=cadmus
REPO?=prologic/$(APP)
TAG?=latest
BUILD?=-dev

all: dev

dev: build
	@./cmd/$(APP)/$(APP) localhost:6667

deps:
	@go get ./...

build: clean deps
	@echo " -> Building $(TAG)$(BUILD)"
	@cd cmd/$(APP) && go build -tags "netgo static_build" -installsuffix netgo \
		-ldflags "-w -X github.com/$(REPO)/$(PACKAGE).GitCommit=$(COMMIT) -X github.com/$(REPO)/$(PACKAGE).Build=$(BUILD)" .
	@echo "Built $$(./cmd/$(APP)/$(APP) -v)"

image:
	@docker build --build-arg TAG=$(TAG) --build-arg BUILD=$(BUILD) -t $(REPO):$(TAG) .
	@echo "Image created: $(REPO):$(TAG)"

test:
	@go test -v -cover -race $(TEST_ARGS)

clean:
	@rm -rf $(APP)
