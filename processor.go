package lambda

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
)

type Response struct {
	Payload Payload `json:"offline_payload"`
}

type Payload struct {
	Error   string      `json:"error,omitempty"`
	Success interface{} `json:"success,omitempty"`
}

func Start(h interface{}) {
	response := Response{}

	authorizer := os.Getenv("IS_LAMBDA_AUTHORIZER")
	requestAuthorizer := os.Getenv("IS_LAMBDA_REQUEST_AUTHORIZER")
	tokenAuthorizer := os.Getenv("IS_LAMBDA_TOKEN_AUTHORIZER")

	// Most common - probably
	if authorizer != "true" {
		response = api(h.(func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)))
	}

	if authorizer == "true" && tokenAuthorizer == "true" {
		response = token(h.(func(request events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error)))
	}

	if authorizer == "true" && requestAuthorizer == "true" {
		response = request(h.(func(request events.APIGatewayCustomAuthorizerRequestTypeRequest) (events.APIGatewayCustomAuthorizerResponse, error)))
	}

	out, _ := json.Marshal(response)
	fmt.Println(string(out))
}
