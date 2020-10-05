.PHONY: build clean deploy

build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/fetch_external fetch_external/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/post_ddb post_ddb/main.go

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
