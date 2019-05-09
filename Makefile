export GO111MODULE=on

test:
	./scripts/codecov.sh

build:
	go build
	cd misc && go build && cd -
	cd chancall && go build && cd -
	cd network && go build && cd -
	cd tool && go build && cd -
	cd test && go build && cd -
