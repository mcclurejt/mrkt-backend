.PHONY: build clean deploy

build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/symbol/company symbol/company/company.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/symbol/historical symbol/historical/historical.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/symbol/stats symbol/stats/stats.go

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
