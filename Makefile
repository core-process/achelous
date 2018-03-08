# package
PACKAGE = github.com/core-process/achelous

# directories
SRC = .build/src/$(PACKAGE)
BIN = .build/bin

# files
SOURCES_GO = $(shell find $(SRC)/ -type f -name '*.go')
SOURCES_C  = $(shell find $(SRC)/ -type f -name '*.c')

BINARIES = \
	$(BIN)/spring-core $(BIN)/spring \
	$(BIN)/upstream-core $(BIN)/upstream

# environment
export GOPATH       = $(CURDIR)/.build
export ARCHITECTURE = $(subst x86_64,amd64,$(shell uname -m))

# settings
include settings.mk
-include settings.user.mk

# build target
all: build

build: $(BINARIES)

## build go sources
$(BIN)/spring-core $(BIN)/upstream-core: $(SRC)/common/config/config.go $(SOURCES_GO) $(SRC)/vendor Makefile | $(BIN)
	cd $(SRC) && $(GO) build -o $@ $(notdir $@)/main.go

$(SRC)/common/config/config.go: $(SRC)/common/config/config.go.tpl $(wildcard config.mk)
	envsubst < $< > $@

## go vendoring
$(SRC)/vendor: $(SRC)/glide.lock
	$(GLIDE) install
	touch $@

$(SRC)/glide.lock: $(SRC)/glide.yaml
	$(GLIDE) update
	touch $@

## build c sources
$(BIN)/spring $(BIN)/upstream: $(SRC)/bootstrap/main.c $(SRC)/bootstrap/config.h $(SOURCES_C) Makefile | $(BIN)
	gcc $< -o $@

$(SRC)/bootstrap/config.h: $(SRC)/bootstrap/config.h.tpl $(wildcard config.mk)
	envsubst < $< > $@

## prepare binary directory
$(BIN):
	mkdir -p $@

# dist target
DEB = .build/dist/achelous_$(VERSION)_$(ARCHITECTURE).deb

dist: $(DEB)

$(DEB): CONTENT = .build/dist/content
$(DEB): $(BINARIES) meta/deb/*
	# assemble files
	mkdir -p $(CONTENT)/usr/sbin
	for bin in $(BINARIES); do \
		cp "$$bin" "$(CONTENT)/usr/sbin/achelous-$$(basename $$bin)"; \
	done
	for alias in sendmail mailq newaliases; do \
		ln -sf "achelous-spring" "$(CONTENT)/usr/sbin/$$alias"; \
	done
	# pack deb
	mkdir -p $(CONTENT)/DEBIAN
	cp meta/deb/* $(CONTENT)/DEBIAN/
	envsubst < meta/deb/control > $(CONTENT)/DEBIAN/control
	cd .build/dist && dpkg-deb --build content $(notdir $@)

# cleanup target
clean:
	rm -rf .build/bin
	rm -rf .build/dist
	rm -rf .build/pkg
	rm -f bootstrap/config.h
	rm -f common/config/config.go
	rm -rf vendor
