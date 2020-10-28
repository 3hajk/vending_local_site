GOENV=GOOS=linux GOARCH=arm GOARM=7
GOCMD=go
GOBUILD=$(GOCMD) build
GOBUILDARM=$(GOENV) CGO_ENABLED=1 CC=arm-linux-gnueabi-gcc $(GOBUILD)
GOLDFLAGS += -linkmode external
GOLDFLAGS += -extldflags
GOLDFLAGS += '-static'
GOFLAGS = -a -ldflags "-s -w $(GOLDFLAGS)"

.PHONY: default
default: build-rpi ;

prep:
	@mkdir -p ./build

setversion:

build:
	$(GOBUILD) -o build/vending_machine app/main.go

build-rpi:
	@echo "Building binary for raspberry pi ..."
	$(GOBUILDARM) $(GOFLAGS) -o build/vending_machine_rpi app/main.go
run:
	./build/vending_machine

test: run

default: prep build

clean:
	@rm -rf ./build