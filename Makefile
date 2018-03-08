# package name
PACKAGE   = github.com/core-process/achelous

# directories
BUILDSPACE = .build
SOURCES    = $(BUILDSPACE)/src/$(PACKAGE)
BINARIES   = $(BUILDSPACE)/bin
PACKSPACE  = .pack

# files
GO_FILES = $(shell [ -d $(SOURCES)/ ] && find $(SOURCES)/ -type f -name '*.go' || true)
C_FILES  = $(shell [ -d $(SOURCES)/ ] && find $(SOURCES)/ -type f -name '*.c' || true)

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

# build target
all: build

build: $(BINARY_FILES)

$(BINARIES)/spring-core: $(GO_FILES) $(SOURCES)/vendor | $(SOURCES) $(BINARIES)
	cd $(SOURCES) && $(GO) build -o $@ spring-core/main.go

$(BINARIES)/upstream-core: $(GO_FILES) $(SOURCES)/vendor | $(SOURCES) $(BINARIES)
	cd $(SOURCES) && $(GO) build -o $@ upstream-core/main.go

$(BINARIES)/spring: $(C_FILES) $(SOURCES)/bootstrap/main.c | $(SOURCES) $(BINARIES)
	gcc $(SOURCES)/bootstrap/main.c -o $@

$(BINARIES)/upstream: $(C_FILES) $(SOURCES)/bootstrap/main.c | $(SOURCES) $(BINARIES)
	gcc $(SOURCES)/bootstrap/main.c -o $@

$(SOURCES)/glide.lock: $(SOURCES)/glide.yaml | $(SOURCES)
	cd $(SOURCES) && $(GLIDE) update
	touch $@

$(SOURCES)/glide.yaml: | $(SOURCES)

$(SOURCES)/vendor: $(SOURCES)/glide.lock | $(SOURCES)
	cd $(SOURCES) && $(GLIDE) install
	@touch $@

$(SOURCES):
	mkdir -p $(dir $@)
	ln -sf $(CURDIR) $@

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
	rm -rf $(BUILDSPACE)
	rm -rf $(PACKSPACE)/.build
