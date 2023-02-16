package streaming

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
)

// const AWS = require("aws-sdk");
// AWS.config.update({ region: process.env.AWS_REGION });

// class StreamerConnection {
//   constructor(domain, stage) {
//     this.apigwManagementApi = new AWS.ApiGatewayManagementApi({
//       apiVersion: "2018-11-29",
//       endpoint: domain + '/' + stage
//     });
//   }

//   async postData(data, connectionId) {
//     return await this.apigwManagementApi.postToConnection({
//       ConnectionId: connectionId,
//       Data: data
//     }).promise();
//   }
// }

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

func (sc *StreamerConnection) PostData(data []byte, connectionId *string) {

	param := apigatewaymanagementapi.PostToConnectionInput{ConnectionId: connectionId, Data: data}

}
