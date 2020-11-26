.EXPORT_ALL_VARIABLES:
GOBIN=${HOME}/go/bin/

deps:
	go get -v

build:
	go build

clean:
	go clean -x -i github.com/danie1sullivan/gooom

install:
	go install
