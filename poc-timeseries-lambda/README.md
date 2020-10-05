## bash commands to run poc lambda

GOOS=linux go build main.go
zip poc-timeseries.zip main
aws lambda update-function-code --function-name poc-timeseries --zip-file fileb://poc-timeseries.zip
