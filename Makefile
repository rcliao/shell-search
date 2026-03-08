.PHONY: build test vet clean

BINARY := shell-search
PKG := ./cmd/shell-search

build:
	go build -o $(BINARY) $(PKG)

test:
	go test ./...

vet:
	go vet ./...

clean:
	rm -f $(BINARY)
