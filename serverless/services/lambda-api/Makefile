.PHONY: build clean deploy

build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/symbols symbols/symbols.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/historical historical/historical.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/company company/company.go
clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
