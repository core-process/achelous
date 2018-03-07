PACKAGE   = github.com/core-process/achelous

WORKSPACE = $(realpath $(dir $(realpath $(firstword $(MAKEFILE_LIST)))))
GOPATH    = $(WORKSPACE)/.gopath
SOURCES   = $(GOPATH)/src/$(PACKAGE)
BINARIES  = $(GOPATH)/bin

GO        = go
GLIDE     = glide

all: build

build: $(BINARIES)/spring-core $(BINARIES)/upstream-core

$(BINARIES)/spring-core: $(SOURCES)/vendor | $(SOURCES) $(BINARIES)
	cd $(SOURCES) && $(GO) build -o $@ spring-core/main.go

$(BINARIES)/upstream-core: $(SOURCES)/vendor | $(SOURCES) $(BINARIES)
	cd $(SOURCES) && $(GO) build -o $@ upstream-core/main.go

$(SOURCES)/glide.lock: $(SOURCES)/glide.yaml | $(SOURCES)
	cd $(SOURCES) && $(GLIDE) update
	touch $@

$(SOURCES)/glide.yaml: | $(SOURCES)

$(SOURCES)/vendor: $(SOURCES)/glide.lock | $(SOURCES)
	cd $(SOURCES) && $(GLIDE) install
	@touch $@

$(SOURCES):
	mkdir -p $(dir $@)
	ln -sf $(WORKSPACE) $@

$(BINARIES):
	mkdir -p $@

clean:
	rm -rf $(GOPATH)
