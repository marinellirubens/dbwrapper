VERSION=$(shell cat ./VERSION)

build:
	@go build -o ./bin/${VERSION}/dbwrapper main.go

run: build
	@echo "Running server using make run steps"
	@if [[ -f config.example.json && ! -f config.json ]]; then \
		cp config.example.json config.json; \
	fi
	@./bin/dbwrapper -f config.json

tests:
	@go test -v $(shell go list ./...)
