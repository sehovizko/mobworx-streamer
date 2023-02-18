package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/go-redis/redis"
	"log"
	"os"
)

type MunitStats struct {
	Data []string `json:"data"`
}

const (
	MunitStatsIDXStart = 0
	MunitStatsIDXEnd   = 1
)

var (
	redisClient *redis.Client
	sess        *session.Session
)

func HandleQueryMunitStats(_ aws.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	mStats, err := redisClient.ZRange("munitstatsidx", MunitStatsIDXStart, MunitStatsIDXEnd).Result()
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, nil
	}

	body, err := json.Marshal(MunitStats{
		Data: mStats,
	})
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Access-Control-Allow-Headers": "Content-Type",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "OPTIONS,GET",
		},
		Body: string(body),
	}, nil
}

func main() {
	sess = session.Must(session.NewSession())
	sm := secretsmanager.New(sess)
	redisCredentials, err := sm.GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId: aws.String(os.Getenv("REDIS_CREDENTIALS")),
	})
	if err != nil {
		panic(err)
	}

	log.Println(redisCredentials.String())

	redisClient = redis.NewClient(&redis.Options{})
	lambda.Start(HandleQueryMunitStats)
}
