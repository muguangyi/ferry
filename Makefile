export GO111MODULE=on

test:
	./scripts/codecov.sh

build:
	cd misc && go build && cd -
	cd chancall && go build && cd -
	cd network && go build && cd -
	cd seek && go build && cd -
	cd tool && go build && cd -
