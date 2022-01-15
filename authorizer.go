package lambda

import (
	"os"

	"github.com/aws/aws-lambda-go/events"
)

func token(h func(request events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error)) Response {
	response := Response{
		Payload: Payload{},
	}

	event := events.APIGatewayCustomAuthorizerRequest{}
	err := decode(os.Getenv("LAMBDA_EVENT"), &event)
	if err != nil {
		response.Payload.Error = err.Error()
		return response
	}

	res, err := h(event)
	if err != nil {
		response.Payload.Error = err.Error()
		return response
	}

	response.Payload.Success = res
	return response
}

func request(h func(request events.APIGatewayCustomAuthorizerRequestTypeRequest) (events.APIGatewayCustomAuthorizerResponse, error)) Response {
	response := Response{
		Payload: Payload{},
	}

	event := events.APIGatewayCustomAuthorizerRequestTypeRequest{}
	err := decode(os.Getenv("LAMBDA_EVENT"), &event)
	if err != nil {
		response.Payload.Error = err.Error()
		return response
	}

	res, err := h(event)
	if err != nil {
		response.Payload.Error = err.Error()
		return response
	}

	response.Payload.Success = res
	return response
}
