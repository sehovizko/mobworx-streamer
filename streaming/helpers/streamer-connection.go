package helpers

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
)

type StreamerConnection struct {
	ApiGwManagementApi *apigatewaymanagementapi.ApiGatewayManagementApi
}

var mySession *session.Session

func InitSession() {
	mySession = session.Must(session.NewSession())
}
func GetStreamerConnection(domain string, stage string) *StreamerConnection {
	if mySession == nil {
		panic("Session must be initialized first.")
	}
	apiGateway := apigatewaymanagementapi.New(mySession, aws.NewConfig().WithRegion("us-west-2"))
	return &StreamerConnection{
		ApiGwManagementApi: apiGateway,
	}
}

func (sc *StreamerConnection) PostData(data *[]byte, connectionId string) (*apigatewaymanagementapi.PostToConnectionOutput, error) {
	postToConnectionInput := &apigatewaymanagementapi.PostToConnectionInput{
		ConnectionId: &connectionId,
		Data:         *data,
	}
	return sc.ApiGwManagementApi.PostToConnection(postToConnectionInput)
}
