PRJ=src/kai-ab

GOENV=
GOCMD=$(GOENV) go
GOBUILD=$(GOCMD) install
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod

BUILD_FLGS=-a -tags netgo -installsuffix netgo -ldflags='-extldflags="static"'

SRCS := $(shell find . -name '*.go' -type f)
BINS := $(shell test -d ./bin && find ./bin/ -type f)

all: test build ## test & build

build: $(SRCS) ## build to linux binary
	cd $(CURDIR)/$(PRJ); GOOS=linux GOARCH=amd64 $(GOBUILD) $(BUILD_FLGS) ./exec/...
	cd $(CURDIR)/$(PRJ); GOOS=linux GOARCH=arm64 $(GOBUILD) $(BUILD_FLGS) ./exec/...

.PHONY: test
test: ## run test
	$(GOTEST) -count=1 ./src/kai-ab/...

.PHONY: clean
clean: $(BINS) ## cleanup
	$(GOCLEAN)
	rm -f $(BINS)

xross-build: build build-windows build-mac ## build to windows and mac and linux
build-windows: ## build to windows binary
	cd $(CURDIR)/$(PRJ); CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) $(BUILD_FLGS) ./exec/...
build-mac: ## build to mac binary
	cd $(CURDIR)/$(PRJ); CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(GOBUILD) $(BUILD_FLGS) ./exec/...

mod: $(CURDIR)/$(PRJ)/go.mod ## mod ensure
	cd $(CURDIR)/$(PRJ); $(GOMOD) tidy
	cd $(CURDIR)/$(PRJ); $(GOMOD) vendor
	cd $(CURDIR)/$(PRJ); $(GOMOD) graph
modinit: ## mod init
	cd $(CURDIR)/$(PRJ); $(GOMOD) init

.PHONY: help
	all: help
help: ## help
	@awk -F ':|##' '/^[^\t].+?:.*?##/ {\
		printf "\033[36m%-30s\033[0m %s\n", $$1, $$NF \
	}' $(MAKEFILE_LIST)
