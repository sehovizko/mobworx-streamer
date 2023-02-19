package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/go-redis/redis/v8"
	"os"
)

type MunitStats struct {
	Data []string `json:"data"`
}

const (
	MunitStatsIDXStart = 0
	MunitStatsIDXEnd   = 1
)

var redisClient *redis.Client

func HandleQueryMunitStats(ctx aws.Context, _ events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	mStats, err := redisClient.ZRange(ctx, "munitstatsidx", MunitStatsIDXStart, MunitStatsIDXEnd).Result()
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
	redisClient = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDRESS"),
	})
	lambda.Start(HandleQueryMunitStats)
}
