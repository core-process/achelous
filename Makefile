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
$(SRC)/glide.lock: $(SRC)/glide.yaml
	cd $(SRC) && $(GLIDE) update
	touch $@

$(SRC)/vendor: $(SRC)/glide.lock
	cd $(SRC) && $(GLIDE) install
	touch $@

## build c sources
$(BIN)/spring $(BIN)/upstream: $(SRC)/bootstrap/main.c $(SRC)/bootstrap/config.h $(SOURCES_C) Makefile | $(BIN)
	gcc $< -o $@

$(SRC)/bootstrap/config.h: $(SRC)/bootstrap/config.h.tpl $(wildcard config.mk)
	envsubst < $< > $@

## prepare binary directory
$(BIN):
	mkdir -p $@

# pack target
pack: .pack/.build/achelous_$(VERSION)_$(ARCHITECTURE).deb

.pack/.build/achelous_$(VERSION)_$(ARCHITECTURE).deb: $(BINARIES) .pack/deb/*
	# assemble files
	mkdir -p .pack/.build/root/usr/sbin
	for bin in $(BINARIES); do \
		cp "$$bin" ".pack/.build/root/usr/sbin/achelous-$$(basename $$bin)"; \
	done
	for alias in sendmail mailq newaliases; do \
		ln -sf "achelous-spring" ".pack/.build/root/usr/sbin/$$alias"; \
	done
	# pack deb
	mkdir -p .pack/.build/root/DEBIAN
	cp .pack/deb/* .pack/.build/root/DEBIAN/
	envsubst < .pack/deb/control > .pack/.build/root/DEBIAN/control
	cd .pack/.build && dpkg-deb --build root achelous_$(VERSION)_$(ARCHITECTURE).deb

# cleanup target
clean:
	rm -rf .build/bin
	rm -rf .pack/.build
	rm -f bootstrap/config.h
	rm -f common/config/config.go
	rm -rf vendor
