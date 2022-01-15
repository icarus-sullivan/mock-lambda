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

type Processor interface {
	Process(h interface{}) Response
}

func Start(h interface{}) {
	response := Response{}
	fmt.Println("IS_LAMBDA_AUTHORIZER", os.Getenv("IS_LAMBDA_AUTHORIZER"))
	// Most common - probably
	if os.Getenv("IS_LAMBDA_AUTHORIZER") != "true" {
		response = Api(h.(func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)))
	}
	if os.Getenv("IS_LAMBDA_AUTHORIZER") == "true" {
		response = Authorizer(h.(func(request events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error)))
	}

	out, _ := json.Marshal(response)
	fmt.Println(string(out))
}
