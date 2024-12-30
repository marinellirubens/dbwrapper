build:
	@go build -o ./bin/dbwrapper cmd/app/main.go

run: build
	@echo "Running server using make run steps"
	@if [[ -f config.example.ini && ! -f config.ini ]]; then \
		cp config.example.ini config.ini; \
	fi
	@./bin/dbwrapper -f config.ini

