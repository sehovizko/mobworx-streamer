package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"log"
)

func HandleQueryAdminRooms(_ aws.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("%+v", event)
	return events.APIGatewayProxyResponse{}, nil
}

func main() {
	lambda.Start(HandleQueryAdminRooms)
}
