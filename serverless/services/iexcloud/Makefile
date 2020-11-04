.PHONY: build clean deploy

build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/iex-symbols iex-symbols/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/batch-single-day batch-single-day/main.go

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
