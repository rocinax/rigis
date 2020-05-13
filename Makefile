# define application
$(eval APPNAME := $(shell basename `pwd`))
$(eval APPVERSION := $(shell cat VERSION))
$(eval VERSIONCHECK := $(shell sh build/check_version_format.sh $(APPVERSION)))
APPFULLNAME=$(APPNAME)-$(APPVERSION)

# Go build commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Develop Environment Directory
CMDDIR=cmd
BUILDWORKDIR=work
RELEASEDIR=release

# Production Environment Directory
BINARYDIR=bin
TEMPLATEDIR=template
CONFIGDIR=config
LOGDIR=logs
INITDIR=init
$(eval CMDS := $(shell find $(CMDDIR)/* -type d | sed 's/$(CMDDIR)\///'))

.PHONY: all
all: check release

# make check ... check build environment
.PHONY: check
check:
	$(if $(VERSIONCHECK),$(error Invalid version format))
	$(GOTEST) -v ./...

# make install ... install application
.PHONY: install
install: clean build
	mkdir -p $(ROOTDIR)/$(APPNAME)
	mkdir -p $(ROOTDIR)/$(APPNAME)
	cp -Rp $(BUILDWORKDIR)/$(APPNAME) /opt/$(APPNAME)/

# make clean ... clean application directory
.PHONY: clean
clean:
	$(GOCLEAN)
	rm -rf $(BUILDWORKDIR)/*

# make build ... build application
.PHONY: build
build: check
	mkdir -p $(BUILDWORKDIR)/$(APPFULLNAME)/$(BINARYDIR)
	mkdir -p $(BUILDWORKDIR)/$(APPFULLNAME)/$(TEMPLATEDIR)
	mkdir -p $(BUILDWORKDIR)/$(APPFULLNAME)/$(CONFIGDIR)
	mkdir -p $(BUILDWORKDIR)/$(APPFULLNAME)/$(LOGDIR)
	mkdir -p $(BUILDWORKDIR)/$(APPFULLNAME)/$(INITDIR)
	cp -Rp $(TEMPLATEDIR) $(BUILDWORKDIR)/$(APPFULLNAME)/
	cp -Rp $(CONFIGDIR) $(BUILDWORKDIR)/$(APPFULLNAME)/
	cp -Rp $(INITDIR) $(BUILDWORKDIR)/$(APPFULLNAME)/
	cp -p VERSION $(BUILDWORKDIR)/$(APPFULLNAME)/
	$(foreach CMD,$(CMDS),CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILDWORKDIR)/$(APPFULLNAME)/$(BINARYDIR)/$(CMD) -v $(CMDDIR)/$(CMD)/main.go)
	
# release source files
.PHONY: release
release: clean build
	rm -f $(RELEASEDIR)/$(APPFULLNAME).tar.gz
	mkdir -p $(RELEASEDIR)
	tar -zcvf $(RELEASEDIR)/$(APPFULLNAME).tar.gz -C $(BUILDWORKDIR)/ .

# cross compile
.PHONY: build-linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARYUNIX) -v