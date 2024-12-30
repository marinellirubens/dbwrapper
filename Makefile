build:
	go build -o ./bin/dbwrapper cmd/app/main.go

run: build
	./bin/dbwrapper -f internal/config/config.ini

