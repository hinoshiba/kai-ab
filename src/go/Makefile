PRJ=src/kai-ab

GOENV=GOPATH=$(CURDIR)
GOCMD=$(GOENV) go
GOBUILD=$(GOCMD) install
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GODEP=$(GOENV) dep

BUILD_FLGS=-a -tags netgo -installsuffix netgo -ldflags='-extldflags="static"'

SRCS := $(shell find . -name '*.go' -type f)
BINS := $(shell test -d ./bin && find ./bin/ -type f)

all: test build ## test & build

build: $(SRCS) ## build to linux binary
	$(GOBUILD) $(BUILD_FLGS) ./src/kai-ab/exec/...

.PHONY: test
test: ## run test
	$(GOTEST) -count=1 ./src/kai-ab/...

.PHONY: clean
clean: $(BINS) ## cleanup
	$(GOCLEAN)
	rm -f $(BINS)

xross-build: build build-windows build-mac ## build to windows and mac and linux
build-windows: ## build to windows binary
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) $(BUILD_FLGS) ./src/kai-ab/exec/...
build-mac: ## build to mac binary
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) $(BUILD_FLGS) ./src/kai-ab/exec/...

dep: $(CURDIR)/$(PRJ)/Gopkg.toml ## dep ensure
	cd $(CURDIR)/$(PRJ); $(GODEP) ensure
	cd $(CURDIR)/$(PRJ); $(GODEP) status
depinit: ## dep init
	cd $(CURDIR)/$(PRJ); $(GODEP) init

.PHONY: help
	all: help
help: ## help
	@awk -F ':|##' '/^[^\t].+?:.*?##/ {\
		printf "\033[36m%-30s\033[0m %s\n", $$1, $$NF \
	}' $(MAKEFILE_LIST)
