GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod

BINARY_NAME :=
ifeq ($(OS),Windows_NT)
	BINARY_NAME += euchred.exe
else
	BINARY_NAME += euchred
endif

.PHONY: all
all: test $(BINARY_NAME)

.PHONY: test
test: deps
	$(GOTEST) -v ./...

$(BINARY_NAME): deps assets
	$(GOBUILD) -o $(BINARY_NAME) -v .

.PHONY: deps
deps:
	$(GOMOD) download

.PHONY: assets
assets:
	pkger

.PHONY: run
run: $(BINARY_NAME)
	./$(BINARY_NAME)
