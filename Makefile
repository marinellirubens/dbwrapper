VERSION=$(shell cat ./VERSION)

build:
	@go build -o ./bin/${VERSION}/dbwrapper main.go

run: build
	@echo "Running server using make run steps"
	@if [[ -f config.example.ini && ! -f config.ini ]]; then \
		cp config.example.ini config.ini; \
	fi
	@./bin/dbwrapper -f config.ini

tests:
	@go test -v $(shell go list ./...)
