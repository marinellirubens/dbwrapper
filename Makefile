VERSION=$(shell cat ./VERSION)

build:
	go build -o ./bin/${VERSION}/dbwrapper main.go
	if [[ -f config.example.json && ! -f config.json ]]; then \
		cp config.example.json config.json; \
	fi

run: build
	echo "Running server using make run steps"
	./bin/dbwrapper -f config.json

tests:
	go test -v $(shell go list ./...)

logbuild:
	cp logrotate /etc/logrotate.d/dbwrapper

build_container: build
	./build_container.sh

create_container: build_container
	./create_container.sh

