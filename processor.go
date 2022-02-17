package lambda

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

type Response struct {
	Payload Payload `json:"offline_payload"`
}

type Payload struct {
	Error   string      `json:"error,omitempty"`
	Success interface{} `json:"success,omitempty"`
}

type MockLambda struct {
	api     func(h func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) Response
	token   func(h func(request events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error)) Response
	request func(h func(request events.APIGatewayCustomAuthorizerRequestTypeRequest) (events.APIGatewayCustomAuthorizerResponse, error)) Response
	sqs     func(h func(ctx context.Context, request events.SQSEvent) error) Response
}

func (ml *MockLambda) start(h interface{}) Response {
	response := Response{}

	types := reflect.TypeOf(h)
	inputCount := types.NumIn()
	inputTypes := make([]string, inputCount)

	for i := 0; i < inputCount; i++ {
		inputTypes[i] = types.In(i).String()
	}

	if inputCount == 2 && inputTypes[0] == "context.Context" && inputTypes[1] == "events.APIGatewayProxyRequest" {
		response = ml.api(h.(func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)))
	} else if inputCount == 1 && inputTypes[0] == "events.APIGatewayCustomAuthorizerRequest" {
		response = ml.token(h.(func(request events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error)))
	} else if inputCount == 1 && inputTypes[0] == "events.APIGatewayCustomAuthorizerRequestTypeRequest" {
		response = ml.request(h.(func(request events.APIGatewayCustomAuthorizerRequestTypeRequest) (events.APIGatewayCustomAuthorizerResponse, error)))
	} else if inputCount == 2 && inputTypes[0] == "context.Context" && inputTypes[1] == "events.SQSEvent" {
		response = ml.sqs(h.(func(ctx context.Context, request events.SQSEvent) error))
	} else {
		response.Payload.Error = "no handler found for method signature func(" + strings.Join(inputTypes, ", ") + ")"
	}

	return response
}

func Start(h interface{}) {
	ml := MockLambda{api: api, token: token, request: request, sqs: sqs}
	response := ml.start(h)

	out, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(string(out))
}
