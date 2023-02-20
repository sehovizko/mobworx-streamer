package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"log"
)

func HandleUploadPart(_ aws.Context, event events.APIGatewayProxyRequest) error {
	log.Printf("connection domain name: %s, stage: %s", event.RequestContext.DomainName, event.RequestContext.Stage)
	return nil
}

func main() {
	lambda.Start(HandleUploadPart)
}
