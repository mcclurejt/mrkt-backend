.PHONY: build clean deploy

build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/publish-historical-for-symbol publish-historical-for-symbol/main.go

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
