GOENV=GOOS=linux GOARCH=arm GOARM=7
GOCMD=go
GOBUILD=$(GOCMD) build
GOBUILDARM=$(GOENV) CGO_ENABLED=1 CC=arm-linux-gnueabi-gcc $(GOBUILD)
GOLDFLAGS += -linkmode external
GOLDFLAGS += -extldflags
GOLDFLAGS += '-static'
VERSION_TAG= $(shell bash ldflags.sh)
GOFLAGS = -a -ldflags "-s -w $(GOLDFLAGS) $(VERSION_TAG)"

.PHONY: default
default: build-rpi ;

version:
	@git describe --tags

prep:
	@mkdir -p ./build

major:
	bash update.sh -v patch
minor:
	bash update.sh -v patch
path:
	bash update.sh -v patch

build:
	$(GOBUILD) -o build/vending_machine app/main.go

raspberry:
	@echo "Building binary for raspberry pi ..."
	$(GOBUILDARM) $(GOFLAGS) -o build/vending_machine_rpi app/main.go

release: raspberry path

run:
	./build/vending_machine

test: run

clean:
	@rm -rf ./build