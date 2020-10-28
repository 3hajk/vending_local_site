GOENV=GOOS=linux GOARCH=arm GOARM=7
GOCMD=go
GOBUILD=$(GOCMD) build
GOBUILDARM=$(GOENV) CGO_ENABLED=1 CC=arm-linux-gnueabi-gcc $(GOBUILD)
GOLDFLAGS += -linkmode external
GOLDFLAGS += -extldflags
GOLDFLAGS += '-static'
VERSION_PKG=main
CURRENT_VERSION=$(shell git describe --abbrev=0 --tags 2>/dev/null)
COMMIT_HASH=$(shell git show -s --format=%h)
COMMIT_STAMP=$(shell git show -s --format=%ct)
BUILD_USER=$(shell id -u -n)
BUILD_HOST=$(shell hostname)
BUILD_STAMP=$(shell date +%s)


GIT_COMMIT = $(shell git rev-parse HEAD)
GIT_SHA    = $(shell git rev-parse --short HEAD)
GIT_TAG    = $(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)
GIT_DIRTY  = $(shell test -n "`git status --porcelain`" && echo "dirty" || echo "clean")


GOFLAGS = -a -ldflags "-s -w $(GOLDFLAGS) $(VFLAGS)"

define vflag
	VFLAGS=$(VFLAGS) -X $(VERSION_PKG).$1=$2
endef

define update
	$(eval v := $(shell git describe --tags --abbrev=0 | sed -Ee 's/^v|-.*//'))
    $(eval n := $(shell echo $(v) | awk -F. -v OFS=. -v f=$1 '{ $$f++ } 1'))
	@echo "Current Version: $(v)"
	@echo "$(2) updating $(v) to $(n)"

	@git tag ${NEW_TAG}
  	@git push --tags
  	@git push
endef


.PHONY: default
default: build-rpi ;

info:
	@echo "Version:           ${VERSION}"
	@echo "Git Tag:           ${GIT_TAG}"
	@echo "Git Commit:        ${GIT_COMMIT}"
	@echo "Git SHA:           ${GIT_SHA}"
	@echo "Git Tree State:    ${GIT_DIRTY}"

version-flag:
	$(eval $(call vflag,Number,$(CURRENT_VERSION)))
    $(eval $(call vflag,CommitHash,$(COMMIT_HASH)))
    $(eval $(call vflag,CommitStamp,$(COMMIT_STAMP)))
    $(eval $(call vflag,BuildUser,$(BUILD_USER)))
    $(eval $(call vflag,BuildHost,$(BUILD_HOST)))
    $(eval $(call vflag,BuildStamp,$(BUILD_STAMP)))


version:
	@git describe --tags

prep:
	@mkdir -p ./build

major:
	$(call update,1,major)
minor:
	$(call update,2,minor)
path:
	$(call update,3,path)
	@echo $(n)

build:
	$(GOBUILD) -o build/vending_machine app/main.go

raspberry:
	@echo "Building binary for raspberry pi ..."
	$(GOBUILDARM) $(GOFLAGS) -o build/vending_machine_rpi app/main.go

release: version-flag raspberry path

next:
	$(eval v := $(shell git describe --tags --abbrev=0 | sed -Ee 's/^v|-.*//'))
ifeq ($1, major)
	$(eval f := 1)
else ifeq ($(bump), minor)
	$(eval f := 2)
else
	$(eval f := 3)
endif
	@echo $(v) | awk -F. -v OFS=. -v f=$(f) '{ $$f++ } 1'

run:
	./build/vending_machine

test: run

clean:
	@rm -rf ./build