
all: clean build
build:
	go build concurry.go
test: 
	go test -v
clean:
	go clean
	rm -f concurry
.PHONY: clean
