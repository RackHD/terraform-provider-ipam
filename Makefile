default: build

deps:
	go get ./...

test: deps
	go test -v ./...

clean:
	rm -f ./bin/terraform-provider-ipam

build: deps
	go build -o bin/terraform-provider-ipam

install: clean build
	cp -f ./bin/terraform-provider-ipam $(shell dirname `which terraform`)
