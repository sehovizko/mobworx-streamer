package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
)

type RenditionType string

const (
	VideoRenditionType RenditionType = "video"
	AudioRenditionType RenditionType = "audio"
)

type MimeType string

const (
	Video       MimeType = "video/mp4"
	Audio       MimeType = "audio/mp4"
	Application MimeType = "application/mp4"
)

var renditionToMimeType = map[RenditionType]MimeType{
	VideoRenditionType: Video,
	AudioRenditionType: Audio,
}

func GetMimeType(renditionType RenditionType) string {
	if value, ok := renditionToMimeType[renditionType]; ok {
		return string(value)
	}
	return string(Application)
}

func HandleUpdateRendition(_ aws.Context, event events.APIGatewayProxyRequest) error {
	return nil
}

func main() {
	lambda.Start(HandleUpdateRendition)
}
