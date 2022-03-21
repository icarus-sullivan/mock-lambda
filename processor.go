package lambda

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"

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
		handler := reflect.ValueOf(h)
		handlerType := reflect.TypeOf(h)
		if handlerType.Kind() != reflect.Func {
			response.Payload.Error = fmt.Sprintf("handler kind %s is not %s", handlerType.Kind(), reflect.Func)
			return response
		}

		takesContext, err := validateArguments(handlerType)
		if err != nil {
			response.Payload.Error = err.Error()
			return response
		}

		if err := validateReturns(handlerType); err != nil {
			response.Payload.Error = err.Error()
			return response
		}

		response = func(ctx context.Context, payload string) Response {
			var (
				args   []reflect.Value
				result Response
			)

			if takesContext {
				args = append(args, reflect.ValueOf(ctx))
			}
			if (handlerType.NumIn() == 1 && !takesContext) || handlerType.NumIn() == 2 {
				eventType := handlerType.In(handlerType.NumIn() - 1)
				event := reflect.New(eventType)

				if err := decode(payload, event.Interface()); err != nil {
					result.Payload.Error = err.Error()
					return result
				}

				args = append(args, event.Elem())
			}

			callResult := handler.Call(args)

			// convert return values into (interface{}, error)
			if len(callResult) > 0 {
				if errVal, ok := callResult[len(callResult)-1].Interface().(error); ok {
					result.Payload.Error = errVal.Error()
				}
			}
			if len(callResult) > 1 {
				result.Payload.Success = callResult[0].Interface()
			}

			return result
		}(context.TODO(), os.Getenv("LAMBDA_EVENT"))
	}

	return response
}

func Start(h interface{}) {
	if h == nil {
		fmt.Println("no handler found")
		return
	}

	ml := MockLambda{api: api, token: token, request: request, sqs: sqs}
	response := ml.start(h)

	out, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(string(out))
}
