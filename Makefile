GOCMD=go
#GOOS=linux
GOOS=darwin
#GOARCH=amd64
GOARCH=arm64
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
BINARY_NAME=app
MAINDOTGO=cmd/web/main.go
PRODDIR=prod

all: help

## help: show this help message
.PHONY: help
help: Makefile
	@echo
	@echo " Choose a make command to run"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

## build: build a binary
.PHONY: build
build: css templ test
	$(GOBUILD) -o $(PRODDIR)/$(BINARY_NAME) -v  $(MAINDOTGO)
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(GOBUILD) -o $(PRODDIR)/$(BINARY_NAME) -v  $(MAINDOTGO)

## test: run go unit tests
.PHONY: test
test: 
	$(GOTEST) -v ./...

## clean: clean the binary and generated css
.PHONY: clean
clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f public/output.css

## run: build and run local project
.PHONY: run
run: build
	cd $(PRODDIR) && ./$(BINARY_NAME)

## css: build tailwindcss
.PHONY: css
css:
	tailwindcss -i css/input.css -o $(PRODDIR)/public/output.css --minify

## templ: templ generate
.PHONY: templ
templ:
	@templ generate
