BUILDPATH=$(CURDIR)
GO=$(shell which go)
GOBUILD=$(GO) build
GOCLEAN=$(GO) clean
GOGET=$(GO) get

EXENAME=main.go

export GOPATH=$(CURDIR)

myname:
	@echo "I am a makefile"

makedir:
	@echo "start building tree..."
	@if [ ! -d $(BUILDPATH)/bin ] ; then mkdir -p $(BUILDPATH)/bin ; fi
	@if [ ! -d $(BUILDPATH)/pkg ] ; then mkdir -p $(BUILDPATH)/pkg ; fi

get:
	@echo "Downloading deps"
	@$(GOGET) github.com/gin-contrib/cors
	@$(GOGET) github.com/gin-gonic/contrib/static
	@$(GOGET) github.com/gin-gonic/gin
	@$(GOGET) github.com/jinzhu/gorm
	@$(GOGET) github.com/mattn/go-sqlite3
	
build:
	@echo "start building..."
	$(GOBUILD) $(EXENAME)
	@echo "Yay! all DONE!"

clean:
	@echo "cleanning"
	@rm -rf $(BUILDPATH)/bin/$(EXENAME)
	@rm -rf $(BUILDPATH)/pkg
	@rm -rf $(BUILDPATH)/src/github.com

all: makedir get build
