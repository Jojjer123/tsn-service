MODULE   = $(shell env GO111MODULE=on $(GO) list -m)
DATE    ?= $(shell date +%FT%T%z)
VERSION ?= $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || \
			cat $(CURDIR)/.version 2> /dev/null || echo v0)
PKGS     = $(or $(PKG),$(shell env GO111MODULE=on $(GO) list ./...))
TESTPKGS = $(shell env GO111MODULE=on $(GO) list -f \
			'{{ if or .TestGoFiles .XTestGoFiles }}{{ .ImportPath }}{{ end }}' \
			$(PKGS))
BIN     = $(CURDIR)/bin
OUT 	= $(CURDIR)/build/_output
MK_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
CUR_DIR := $(notdir $(patsubst %/,%,$(dir $(MK_PATH))))

GO      = go
TIMEOUT = 15
V = 0
Q = $(if $(filter 1,$V),,@)
M = $(shell printf "\033[34;1m▶\033[0m")

# Export environment variables if a .env file is present.
ifeq ($(ENV_EXPORTED),) # ENV vars not yet exported
ifneq ("$(wildcard .env)","")
sinclude .env
export $(shell [ -f .env ] && sed 's/=.*//' .env)
export ENV_EXPORTED=true
$(info Note — An .env file exists. Its contents have been exported as environment variables.)
endif
endif

# defaults if not set in the .env file

GO111MODULE?=on
CGO_ENABLED?=0
PRJ_NAME?=$(CUR_DIR)
PRJ_VERSION?=latest
DOCKERFILE?=./build/$(PRJ_NAME)/Dockerfile

.DEFAULT_GOAL := all

$(BIN):
	@mkdir -p $@

.PHONY: docker-$(PRJ_NAME)
docker-$(PRJ_NAME): ; $(info $(M) building docker image...) @ ## Common build docker image
	docker build . -f $(DOCKERFILE) \
		--build-arg o=./bin/$(basename $(MODULE)) \
		-t $(PRJ_NAME):$(PRJ_VERSION)



# this and the common clean will both executed because of ::
.PHONY: clean
clean:: ; $(info $(M) tsn-service clean) @ ## clean (ADDITIONAL)
	@rm -rf $(BIN)
	@rm -rf $(OUT)

.PHONY: images
images: docker-$(PRJ_NAME) ; $(info $(M) building images...) @ ## build all docker images (ADDITIONAL)

.PHONY: images-push
images-push: images $(DOCKER_LOGIN) ; $(info $(M) pushing images...) @ ## push docker images (PROJECT)
	docker push $(PRJ_NAME):$(PRJ_VERSION)

.PHONY: kind
kind: images ; $(info $(M) add images to kind cluster...) @ ## add images to kind (ADDITIONAL)
	@if [ "`kind get clusters`" = '' ]; then echo "no kind cluster found" && exit 1; fi
	kind load docker-image $(PRJ_NAME):$(PRJ_VERSION)

.PHONY: build
build: $(BIN) ; $(info $(M) building executable...) @ ## Common build program binary
	$Q CGO_ENABLED=$(CGO_ENABLED) $(GO) build \
		-tags release \
		-ldflags '-X $(MODULE)/cmd.version=$(PRJ_VERSION) -X $(MODULE)/cmd.commit=$(VERSION) -X $(MODULE)/cmd.date=$(DATE)' \
		-o $(BIN)/$(basename $(MODULE)) main.go


.PHONY: deploy
deploy: build images images-push kind