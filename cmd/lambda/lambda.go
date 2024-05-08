package main

import (
	"context"
	"encoding/json"
	"fmt"

	"echo-server/cmd/api"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	echoAdapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"
)

var echoLambda *echoAdapter.EchoLambda

func init() {
	app := api.New()
	echoLambda = echoAdapter.New(app)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	jsonData, err := json.Marshal(req)
	if err == nil {
		fmt.Println(string(jsonData))
	}
	return echoLambda.Proxy(req)
}

func main() {
	lambda.Start(Handler)
}
