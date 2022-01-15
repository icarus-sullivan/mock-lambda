package lambda

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
)

type AuthorizerLambda func(request events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error)

type AuthorizerResponse struct {
	AuthorizerPayload AuthorizerPayload `json:"offline_payload"`
}

type AuthorizerPayload struct {
	Error   string                                    `json:"error,omitempty"`
	Success events.APIGatewayCustomAuthorizerResponse `json:"success,omitempty"`
}

func AuthorizerHandler(h AuthorizerLambda) {
	response := AuthorizerResponse{
		AuthorizerPayload: AuthorizerPayload{},
	}

	sanitizedEvent := SanitizeJSON(os.Getenv("LAMBDA_EVENT"))

	event := events.APIGatewayCustomAuthorizerRequest{}
	err := json.Unmarshal([]byte(sanitizedEvent), &event)
	if err != nil {
		response.AuthorizerPayload.Error = err.Error()
		out, _ := json.Marshal(response)
		fmt.Println(string(out))
		os.Exit(1)
	}

	res, err := h(event)
	if err != nil {
		response.AuthorizerPayload.Error = err.Error()
		out, _ := json.Marshal(response)
		fmt.Println(string(out))
		os.Exit(1)
	}

	response.AuthorizerPayload.Success = res
	out, _ := json.Marshal(response)
	fmt.Println(string(out))
}
