.PHONY: build clean deploy

build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/post_ticker_sqs post_ticker_sqs/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/fetch_ticker fetch_ticker/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/post_ddb post_ddb/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/initialize_step_function initialize_step_function/main.go

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
