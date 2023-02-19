package helpers

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
)

type StreamerConnection struct {
	ApiGwManagementApi *apigatewaymanagementapi.ApiGatewayManagementApi
	Session            *session.Session
}

func NewStreamerConnection(session *session.Session) *StreamerConnection {
	apiGateway := apigatewaymanagementapi.New(session, aws.NewConfig().WithRegion("us-west-2"))
	return &StreamerConnection{
		ApiGwManagementApi: apiGateway,
		Session:            session,
	}
}

func (sc *StreamerConnection) PostData(data []byte, connectionId string) (*apigatewaymanagementapi.PostToConnectionOutput, error) {
	postToConnectionInput := &apigatewaymanagementapi.PostToConnectionInput{
		ConnectionId: &connectionId,
		Data:         data,
	}
	return sc.ApiGwManagementApi.PostToConnection(postToConnectionInput)
}
