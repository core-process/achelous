# package
PACKAGE = github.com/core-process/achelous

# directories
SRC = .build/src/$(PACKAGE)
BIN = .build/bin

# binaries
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
$(BIN)/spring-core $(BIN)/upstream-core: $(SRC)/common/config/config.go $(shell find $(SRC)/ -type f -name '*.go') $(SRC)/vendor | $(BIN)
	cd $(SRC) && $(GO) build -buildmode=pie -o $@ $(notdir $@)/main.go

$(SRC)/common/config/config.go: $(SRC)/common/config/config.go.tpl $(wildcard settings*.mk)
	envsubst < $< > $@

## go vendoring
$(SRC)/vendor: $(SRC)/glide.lock
	$(GLIDE) install
	touch $@

$(SRC)/glide.lock: $(SRC)/glide.yaml
	$(GLIDE) update
	touch $@

## build c sources
$(BIN)/spring: $(SRC)/bootstrap/spring.c $(SRC)/bootstrap/switchuser.c $(SRC)/bootstrap/coreprocess.c $(SRC)/bootstrap/config.h | $(BIN)
	gcc $(filter %.c,$^) -o $@

$(BIN)/upstream: $(SRC)/bootstrap/upstream.c $(SRC)/bootstrap/switchuser.c $(SRC)/bootstrap/coreprocess.c $(SRC)/bootstrap/daemonise.c $(SRC)/bootstrap/readpid.c $(SRC)/bootstrap/config.h | $(BIN)
	gcc $(filter %.c,$^) -lbsd -o $@

$(SRC)/bootstrap/config.h: $(SRC)/bootstrap/config.h.tpl $(wildcard settings*.mk)
	envsubst < $< > $@

## prepare binary directory
$(BIN):
	mkdir -p $@

# dist target
DEB = .build/dist/achelous_$(VERSION)_$(ARCHITECTURE).deb

# NOTE: use the following command to verify quality of deb file:
# lintian --no-tag-display-limit .build/dist/achelous_1.0-1_amd64.deb

dist: $(DEB)

$(DEB): CONTENT = .build/dist/content
$(DEB): $(BINARIES) meta/deb/*
	# create directories
	mkdir -p $(CONTENT)/DEBIAN
	mkdir -p $(CONTENT)/lib/systemd/system/
	mkdir -p $(CONTENT)/usr/sbin
	mkdir -p $(CONTENT)/usr/share/doc/achelous
	chmod -R 755 $(CONTENT)
	# assemble meta
	envsubst < meta/deb/control > $(CONTENT)/DEBIAN/control
	cp meta/deb/postinst $(CONTENT)/DEBIAN/
	gzip --best -n < meta/deb/changelog > $(CONTENT)/usr/share/doc/achelous/changelog.Debian.gz
	cp meta/deb/copyright $(CONTENT)/usr/share/doc/achelous/
	chmod 644 $(CONTENT)/usr/share/doc/achelous/*
	cp meta/achelous-upstream.service $(CONTENT)/lib/systemd/system/
	chmod 644 $(CONTENT)/lib/systemd/system/*
	# assemble content
	for bin in $(BINARIES); do \
		cp "$$bin" "$(CONTENT)/usr/sbin/achelous-$$(basename $$bin)"; \
	done
	strip $(CONTENT)/usr/sbin/achelous-*
	chmod 755 $(CONTENT)/usr/sbin/achelous-*
	chmod ug+s $(CONTENT)/usr/sbin/achelous-spring
	for alias in sendmail mailq newaliases; do \
		ln -sf "achelous-spring" "$(CONTENT)/usr/sbin/$$alias"; \
	done
	# pack deb
	cd .build/dist && fakeroot dpkg-deb --build content $(notdir $@)

# cleanup target
clean:
	rm -rf .build/bin
	rm -rf .build/dist
	rm -rf .build/pkg
	rm -f bootstrap/config.h
	rm -f common/config/config.go
	rm -rf vendor
	rm -rf testing/testservice/node_modules
	rm -f testing/*.log
