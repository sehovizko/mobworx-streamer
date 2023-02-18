package helpers

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"log"
	"net/http"
	"net/url"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
)

var (
	streamerConn *StreamerConnection
)

type StreamerConnection struct {
	apigwManagementApi *apigatewaymanagementapi.ApiGatewayManagementApi
}

func InitStreamer(domain string, stage string) {
	mySession := session.Must(session.NewSession())
	apiGateway := apigatewaymanagementapi.New(mySession, aws.NewConfig().WithRegion("us-west-2"))
	streamerConn.apigwManagementApi = apiGateway
}

func handleRequest(ctx context.Context, req *events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {

	log.Printf("Handling API Gateway Websocket Proxy Request: %+v", req)
	log.Println()

	// Extract the request information:
	connectionID := req.RequestContext.ConnectionID
	callbackURL := url.URL{
		Scheme: "https",
		Host:   req.RequestContext.DomainName,
		Path:   req.RequestContext.Stage,
	}

	log.Println("Creating API Gateway client for callback URL: ", callbackURL.String())
	mySession := session.Must(session.NewSession())
	apiClient := apigatewaymanagementapi.New(mySession, aws.NewConfig().WithRegion("us-west-2"))

	log.Printf("Created API Gateway Client: %+v", apiClient)
	log.Println()

	// Post a message to the connection:
	_, err := apiClient.PostToConnection(&apigatewaymanagementapi.PostToConnectionInput{
		ConnectionId: aws.String(connectionID),
		Data:         []byte("Test Post to Connection"),
	})
	if err != nil {
		log.Println("Short circuiting and returning 500 because failed to post to connection with error: ", err.Error())
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}
	log.Println("Posted test message to connection")

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "",
	}, nil
}
