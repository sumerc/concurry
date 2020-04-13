
BINARY=concurry
ifeq ($(GOOS), windows)
	ARCH_EXT = .exe
endif

PLATFORMS=darwin linux windows
ARCHITECTURES=386 amd64

#$(info $(GOOS) "windows")

all: clean build
build:
	go build -o $(BINARY)$(ARCH_EXT) -v
build_all:
	$(foreach GOOS, $(PLATFORMS),\
	$(foreach GOARCH, $(ARCHITECTURES), $(env GOOS=$(GOOS); env=$(GOARCH); go build -v -o dist/$(BINARY)-$(GOOS)-$(GOARCH)$(ARCH_EXT))))
test: 
	go test -v
clean:
	go clean
	rm -f concurry
.PHONY: clean

install:
	go install
.PHONY: install
