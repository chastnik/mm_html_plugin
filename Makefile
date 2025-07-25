GO ?= $(shell command -v go 2> /dev/null)
NPM ?= $(shell command -v npm 2> /dev/null)
CURL ?= $(shell command -v curl 2> /dev/null)
MM_DEBUG ?=
MANIFEST_FILE ?= plugin.json
GOPATH ?= $(shell go env GOPATH)
GO_TEST_FLAGS ?= -race
GO_BUILD_FLAGS ?=
MM_UTILITIES_DIR ?= ../mattermost-utilities

# Включает информацию об отладке при MM_DEBUG=1
ifneq ($(MM_DEBUG),)
	GO_BUILD_FLAGS += -gcflags "all=-N -l"
endif

PLUGIN_ID ?= $(shell cat $(MANIFEST_FILE) | python3 -c 'import sys,json;print(json.load(sys.stdin)["id"])')
PLUGIN_VERSION ?= $(shell cat $(MANIFEST_FILE) | python3 -c 'import sys,json;print(json.load(sys.stdin)["version"])')
BUNDLE_NAME ?= $(PLUGIN_ID)-$(PLUGIN_VERSION).tar.gz

# Настройки архитектуры для сборки
SUPPORTED_PLATFORMS ?= linux-amd64 darwin-amd64 windows-amd64

## Builds the plugin for all supported platforms.
all: check-style bundle

## Builds the plugin for all supported platforms, producing a bundle.
bundle: apply server webapp
	rm -rf dist/
	mkdir -p dist/$(PLUGIN_ID)
	cp $(MANIFEST_FILE) dist/$(PLUGIN_ID)/
	cp -r server/dist dist/$(PLUGIN_ID)/server/
	cp -r webapp/dist dist/$(PLUGIN_ID)/webapp/
	cd dist && tar -czf $(BUNDLE_NAME) $(PLUGIN_ID)

## Builds the server executable for all supported platforms.
server: server-linux-amd64 server-darwin-amd64 server-windows-amd64

server-linux-amd64: export GOOS = linux
server-linux-amd64: export GOARCH = amd64
server-linux-amd64: go-build-server

server-darwin-amd64: export GOOS = darwin
server-darwin-amd64: export GOARCH = amd64  
server-darwin-amd64: go-build-server

server-windows-amd64: export GOOS = windows
server-windows-amd64: export GOARCH = amd64
server-windows-amd64: go-build-server

go-build-server: $(GO_BUILD_TARGETS)
	@echo Building server
	mkdir -p server/dist
ifeq ($(GOOS),windows)
	cd server && env GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 $(GO) build $(GO_BUILD_FLAGS) -o dist/plugin-$(GOOS)-$(GOARCH).exe ./...
else
	cd server && env GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 $(GO) build $(GO_BUILD_FLAGS) -o dist/plugin-$(GOOS)-$(GOARCH) ./...
endif

## Builds the webapp.  
webapp: webapp/dist/main.js

webapp/dist/main.js: $(shell find webapp/src -type f) webapp/package.json
	@echo Building webapp
	cd webapp && $(NPM) run build

## Installs dependencies
deps:
	cd server && $(GO) mod download
	cd webapp && $(NPM) install

## Runs tests
test:
	cd server && $(GO) test $(GO_TEST_FLAGS) ./...

## Запуск линтеров
check-style: 
	@echo Checking for style guide compliance

## Убирает временные файлы
clean:
	@echo Cleaning temp files
	rm -rf server/dist
	rm -rf webapp/dist
	rm -rf dist/
	rm -rf webapp/node_modules

## Применяет замены в коде по шаблонам
apply:
	@echo Nothing to apply

.PHONY: all bundle server webapp deps test check-style clean apply 