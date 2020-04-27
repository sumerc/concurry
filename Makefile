PLATFORMS=darwin linux windows
ARCHITECTURES=386 amd64
BINARY=concurry

all: clean build
build:
	go build -v -o $(BINARY)
build_all:
	$(foreach GOOS, $(PLATFORMS),\
	$(foreach GOARCH, $(ARCHITECTURES), $(shell export GOOS=$(GOOS); export GOARCH=$(GOARCH); go build -v -o dist/$(BINARY)-$(GOOS)-$(GOARCH))))
test: 
	go test -v
clean:
	go clean
	rm -f concurry
.PHONY: clean
