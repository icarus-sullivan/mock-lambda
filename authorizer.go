package lambda

import (
	"encoding/json"
	"os"

	"github.com/aws/aws-lambda-go/events"
)

func Authorizer(h func(request events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error)) Response {
	response := Response{
		Payload: Payload{},
	}

	sanitizedEvent := SanitizeJSON(os.Getenv("LAMBDA_EVENT"))

	event := events.APIGatewayCustomAuthorizerRequest{}
	err := json.Unmarshal([]byte(sanitizedEvent), &event)
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
