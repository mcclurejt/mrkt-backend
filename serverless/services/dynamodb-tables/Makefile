.PHONY: build clean deploy

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
