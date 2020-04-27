PLATFORMS=darwin linux windows
ARCHITECTURES=386 amd64

all: clean build
build:
	go build concurry.go
build_all:
	$(foreach GOOS, $(PLATFORMS),\
	$(foreach GOARCH, $(ARCHITECTURES), $(env GOOS=$(GOOS); env=$(GOARCH); go build -v -o dist/$(BINARY)-$(GOOS)-$(GOARCH))))
test: 
	go test -v
clean:
	go clean
	rm -f concurry
.PHONY: clean
