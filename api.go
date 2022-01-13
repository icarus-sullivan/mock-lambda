package lambda

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

type ApiHandler func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

type Response struct {
	Error   string                         `json:"error,omitempty"`
	Success events.APIGatewayProxyResponse `json:"success,omitempty"`
}

func Start(h ApiHandler) {
	response := Response{}

	envEvent := os.Getenv("LAMBDA_EVENT")

	// JSON.stringify escapes this data so we want to unescape to a single backslash
	sanitizedEvent := strings.ReplaceAll(envEvent, "\\\\", "\\")

	fmt.Println("envEvent", sanitizedEvent)
	var event events.APIGatewayProxyRequest
	err := json.Unmarshal([]byte(sanitizedEvent), &event)
	if err != nil {
		response.Error = err.Error()
		out, _ := json.Marshal(response)
		fmt.Println(string(out))
		os.Exit(1)
	}

	res, err := h(context.Background(), event)
	if err != nil {
		response.Error = err.Error()
		out, _ := json.Marshal(response)
		fmt.Println(string(out))
		os.Exit(1)
	}

	response.Success = res
	out, _ := json.Marshal(response)
	fmt.Println(string(out))
}
