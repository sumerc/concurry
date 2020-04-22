
BINARY=concurry

all: clean build
build:
	go build -o $(BINARY) -v
test: 
	go test -v
clean:
	go clean
	rm -f concurry
.PHONY: clean

install:
	go install
.PHONY: install
