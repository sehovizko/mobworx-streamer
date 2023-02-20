package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"hls.streaming.com/src/internal/signals"
	"log"
)

func HandleUploadPart(_ aws.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("connection domain name: %s, stage: %s", event.RequestContext.DomainName, event.RequestContext.Stage)
	message, err := signals.NewDataMessage(event.Body, event.IsBase64Encoded)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	uploadLatency, err := message.UploadLatencyFromNow()
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	log.Println("upload time is ", uploadLatency)

	if !message.Payload.Part.Gap {
		// TODO: call redis client to cache the data with Part cache key
	}

	if message.Payload.Segment.Map != nil {
		// TODO: call redis client to cache the data with Variant or Rendition initial cache key
	}

	// TODO: use redlock to lock by the Variant or Rendition cache key
	// TODO: call updatePart()
	// TODO: unlock redlock by the Variant or Rendition cache key

	// TODO: ack() to the client

	return events.APIGatewayProxyResponse{}, nil
}

func main() {
	lambda.Start(HandleUploadPart)
}
