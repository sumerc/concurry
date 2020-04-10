BINARY_NAME=concurry

all: clean build
build: 
	go -o $(BINARY_NAME) -v
test: 
	go test -v
clean:
	go clean
	rm -f concurry
.PHONY: clean

install:
	cp $(BINARY_NAME) /usr/bin/$(BINARY_NAME)
.PHONY: install
