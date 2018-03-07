# package name
PACKAGE   = github.com/core-process/achelous

# directories
BUILDSPACE = .build
GOPATH     = $(CURDIR)/$(BUILDSPACE)
SOURCES    = $(BUILDSPACE)/src/$(PACKAGE)
BINARIES   = $(BUILDSPACE)/bin

# build tools
GO        = go
GLIDE     = glide
GCC       = gcc

# build target
all: build

build: $(BINARIES)/spring-core $(BINARIES)/upstream-core $(BINARIES)/spring $(BINARIES)/upstream

$(BINARIES)/spring-core: $(SOURCES)/vendor | $(SOURCES) $(BINARIES)
	cd $(SOURCES) && $(GO) build -o $@ spring-core/main.go

$(BINARIES)/upstream-core: $(SOURCES)/vendor | $(SOURCES) $(BINARIES)
	cd $(SOURCES) && $(GO) build -o $@ upstream-core/main.go

$(BINARIES)/spring: | $(SOURCES) $(BINARIES)
	gcc $(SOURCES)/wrapper/main.c -o $@

$(BINARIES)/upstream: | $(SOURCES) $(BINARIES)
	gcc $(SOURCES)/wrapper/main.c -o $@

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

# cleanup target
clean:
	rm -rf $(BUILDSPACE)
