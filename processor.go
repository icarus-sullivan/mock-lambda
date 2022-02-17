package lambda

import (
	"context"
	"encoding/json"
	"errors"
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

func (ml *MockLambda) start(h interface{}) error {
	response := Response{}

	types := reflect.TypeOf(h)
	inputCount := types.NumIn()

	inputTypes := make([]string, inputCount)

	for i := 0; i < inputCount; i++ {
		inputTypes[i] = types.In(i).String()
	}

	fmt.Println(inputTypes)

	if len(inputTypes) == 2 && inputTypes[0] == "context.Context" && inputTypes[1] == "events.APIGatewayProxyRequest" {
		response = ml.api(h.(func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)))
	} else if len(inputTypes) == 1 && inputTypes[0] == "events.APIGatewayCustomAuthorizerRequest" {
		response = ml.token(h.(func(request events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error)))
	} else if len(inputTypes) == 1 && inputTypes[0] == "events.APIGatewayCustomAuthorizerRequestTypeRequest" {
		response = ml.request(h.(func(request events.APIGatewayCustomAuthorizerRequestTypeRequest) (events.APIGatewayCustomAuthorizerResponse, error)))
	} else if len(inputTypes) == 2 && inputTypes[0] == "context.Context" && inputTypes[1] == "events.SQSEvent" {
		response = ml.sqs(h.(func(ctx context.Context, request events.SQSEvent) error))
	} else {
		return errors.New("no handler found for method signature " + strings.Join(inputTypes, ", "))
	}

	out, _ := json.Marshal(response)
	fmt.Println(string(out))

	return nil
}

func Start(h interface{}) {
	ml := MockLambda{api: api, token: token, request: request, sqs: sqs}
	err := ml.start(h)

	if err != nil {
		fmt.Errorf(err.Error())
	}
}
