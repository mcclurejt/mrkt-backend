.PHONY: build clean deploy

build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/initialize_cron initialize_cron/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/fetch_client_data fetch_client_data/main.go

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
