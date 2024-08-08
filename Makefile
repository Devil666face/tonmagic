.DEFAULT_GOAL := help
PROJECT_BIN = $(shell pwd)/bin
PATH := $(PROJECT_BIN):$(PATH)
$(shell [ -f bin ] || mkdir -p $(PROJECT_BIN))
GOBIN = go
GOOS = linux
GOARCH = amd64
LDFLAGS = -extldflags '-static' -w -s -buildid=
GCFLAGS = "all=-trimpath=$(shell pwd) -dwarf=false -l"
ASMFLAGS = "all=-trimpath=$(shell pwd)"
APP := $(PROJECT_BIN)/tonmagic

release: .build .strip ## Build release
debug: .build .strip .copy
build: .build



cert: ## Make ssl cert's
	openssl req -newkey rsa:2048 -nodes -keyout server.key -out server.csr -subj "/CN=your.hostname.here"
	echo "subjectAltName = IP:0.0.0.0" > extfile.cnf
	openssl x509 -req -sha256 -days 365 -in server.csr -signkey server.key -out server.crt -extfile extfile.cnf
	cp server.crt cmd/tonmagic
	cp server.key cmd/tonmagic
	rm \
		extfile.cnf \
		server.csr

help:
	@cat $(MAKEFILE_LIST) | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.build:
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) \
	  $(GOBIN) build -ldflags="$(LDFLAGS)" -trimpath -gcflags=$(GCFLAGS) -asmflags=$(ASMFLAGS) \
	  -o $(APP) cmd/tonmagic/main.go

.strip:
	strip $(APP)
	objcopy --strip-unneeded $(APP)
	upx $(APP)

.copy:
	scp $(APP) host01.d6f.ru:~

