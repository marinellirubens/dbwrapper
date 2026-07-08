VERSION=$(shell cat ./VERSION)

build:
	if [[ -f config.example.json && ! -f src/config.json ]]; then \
		cp config.example.json src/config.json; \
	fi
	cd src && \
	go build -o ./bin/${VERSION}/dbwrapper .

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

