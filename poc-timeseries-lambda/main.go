package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gocarina/gocsv"
)

type Timeseries struct {
	Data []TimeseriesEntry
}
type TimeseriesEntry struct {
	Ticker    string
	Timestamp string `csv:"timestamp"`
	Open      string `csv:"open"`
	Close     string `csv:"close"`
	High      string `csv:"high"`
	Low       string `csv:"low"`
}

func downloadS3File(bucket string, s3File string) {
	//the only writable directory in the lambda is /tmp
	file, err := os.Create("/tmp/" + s3File)
	if err != nil {
		exitErrorf("Unable to open file %q, %v", s3File, err)
	}

	defer file.Close()

	// replace with your bucket region
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	downloader := s3manager.NewDownloader(sess)

	_, err = downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(s3File),
		})
	if err != nil {
		exitErrorf("Unable to download s3File %q, %v", s3File, err)
	}

	fmt.Printf("Downloaded file %q", s3File)
}

func parseCsv(file string) []TimeseriesEntry {
	in, err := os.Open("/tmp/" + file)
	if err != nil {
		panic(err)
	}
	defer in.Close()

	ts := []*TimeseriesEntry{}

	if err := gocsv.UnmarshalFile(in, &ts); err != nil {
		panic(err)
	}
	res := make([]TimeseriesEntry, len(ts))
	for i, row := range ts {
		res[i] = *row
		res[i].Ticker = "AAPL"
	}
	fmt.Printf("%v", res)
	return res
}

func insertIntoDynamoDB(data []TimeseriesEntry) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)

	for _, v := range data {
		fmt.Printf("Timestamp: %v", v.Timestamp)
		fmt.Printf("Ticker: %v", v.Ticker)

		av, err := dynamodbattribute.MarshalMap(v)
		if err != nil {
			exitErrorf("Got error marshalling new item:", av, err)
		}

		tableName := "poc-timeseries"
		input := &dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String(tableName),
		}

		_, err = svc.PutItem(input)
		if err != nil {
			exitErrorf("Got error calling PutItem:", err)
		}
	}
}

func handler(ctx context.Context, s3Event events.S3Event) {
	for _, record := range s3Event.Records {
		s3 := record.S3
		fmt.Printf("[%s - %s] Bucket = %s, Key = %s \n", record.EventSource, record.EventTime, s3.Bucket.Name, s3.Object.Key)
		downloadS3File(s3.Bucket.Name, s3.Object.Key)

		res := parseCsv(s3.Object.Key)
		insertIntoDynamoDB(res)
	}
}

func main() {
	lambda.Start(handler)
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
