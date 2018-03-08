# package name
PACKAGE   = github.com/core-process/achelous

# directories
BUILDSPACE = .build
SOURCES    = $(BUILDSPACE)/src/$(PACKAGE)
BINARIES   = $(BUILDSPACE)/bin
PACKSPACE  = .pack

# files
GO_FILES = $(shell find $(SOURCES)/ -type f -name '*.go')
C_FILES  = $(shell find $(SOURCES)/ -type f -name '*.c')

BINARY_FILES = \
	$(BINARIES)/spring-core $(BINARIES)/spring \
	$(BINARIES)/upstream-core $(BINARIES)/upstream

# build tools
GO        = go
GLIDE     = glide
GCC       = gcc

# environment
export GOPATH       = $(CURDIR)/$(BUILDSPACE)
export ARCHITECTURE = $(subst x86_64,amd64,$(shell uname -m))
export VERSION      = 1.0-1

# configuration
export CONFIG_USER  = achelous
export CONFIG_GROUP = achelous
export CONFIG_SPOOL = /var/spool/achelous
-include config.mk

# build target
all: build

build: $(BINARY_FILES)

## build go binaries
$(BINARIES)/spring-core $(BINARIES)/upstream-core: $(SOURCES)/common/config/config.go $(GO_FILES) $(SOURCES)/vendor | $(BINARIES)
	cd $(SOURCES) && $(GO) build -o $@ $(notdir $@)/main.go

$(SOURCES)/common/config/config.go: $(SOURCES)/common/config/config.go.tpl
	envsubst < $< > $@

## prepare go vendoring
$(SOURCES)/glide.lock: $(SOURCES)/glide.yaml
	cd $(SOURCES) && $(GLIDE) update
	touch $@

$(SOURCES)/vendor: $(SOURCES)/glide.lock
	cd $(SOURCES) && $(GLIDE) install
	touch $@

## build c binaries
$(BINARIES)/spring $(BINARIES)/upstream: $(SOURCES)/bootstrap/main.c $(SOURCES)/bootstrap/config.h $(C_FILES) | $(BINARIES)
	gcc $< -o $@

$(SOURCES)/bootstrap/config.h: $(SOURCES)/bootstrap/config.h.tpl
	envsubst < $< > $@

## prepare directories
$(BINARIES):
	mkdir -p $@

# pack target
pack: $(PACKSPACE)/.build/achelous_$(VERSION)_$(ARCHITECTURE).deb

$(PACKSPACE)/.build/achelous_$(VERSION)_$(ARCHITECTURE).deb: $(BINARY_FILES) $(PACKSPACE)/deb/*
	# assemble files
	mkdir -p $(PACKSPACE)/.build/root/usr/sbin
	cp $(BINARIES)/spring $(PACKSPACE)/.build/root/usr/sbin/achelous-spring
	cp $(BINARIES)/spring-core $(PACKSPACE)/.build/root/usr/sbin/achelous-spring-core
	cp $(BINARIES)/upstream $(PACKSPACE)/.build/root/usr/sbin/achelous-upstream
	cp $(BINARIES)/upstream-core $(PACKSPACE)/.build/root/usr/sbin/achelous-upstream-core
	ln -sf achelous-spring $(PACKSPACE)/.build/root/usr/sbin/sendmail
	ln -sf achelous-spring $(PACKSPACE)/.build/root/usr/sbin/mailq
	ln -sf achelous-spring $(PACKSPACE)/.build/root/usr/sbin/newaliases
	# pack deb
	mkdir -p $(PACKSPACE)/.build/root/DEBIAN
	envsubst < $(PACKSPACE)/deb/control > $(PACKSPACE)/.build/root/DEBIAN/control
	cp $(PACKSPACE)/deb/postinst $(PACKSPACE)/.build/root/DEBIAN/postinst
	cd $(PACKSPACE)/.build && dpkg-deb --build root achelous_$(VERSION)_$(ARCHITECTURE).deb

# cleanup target
clean:
	rm -rf $(BUILDSPACE)/bin
	rm -rf $(PACKSPACE)/.build
	rm -f bootstrap/config.h
	rm -f common/config/config.go
	rm -rf vendor
