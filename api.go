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
	Payload Payload `json:"offline_payload"`
}

type Payload struct {
	Error   string                         `json:"error,omitempty"`
	Success events.APIGatewayProxyResponse `json:"success,omitempty"`
}

func parse(json string) string {
	// When we get this env value it's escaped in a way javascript understands but not
	// go
	// 	- unescape to a single backslash
	return strings.ReplaceAll(json, "\\\\", "\\")
}

func Start(h ApiHandler) {
	response := Response{
		Payload: Payload{},
	}

	envEvent := os.Getenv("LAMBDA_EVENT")

	// JSON.stringify escapes this data so we want to unescape to a single backslash
	sanitizedEvent := parse(envEvent)

	fmt.Println("envEvent", sanitizedEvent)
	var event events.APIGatewayProxyRequest
	err := json.Unmarshal([]byte(sanitizedEvent), &event)
	if err != nil {
		response.Payload.Error = err.Error()
		out, _ := json.Marshal(response)
		fmt.Println(string(out))
		os.Exit(1)
	}
	// {"awsRequestId":"ckydzwmfh0002impyhfafadvi","callbackWaitsForEmptyEventLoop":true,"clientContext":null,"functionName":"oauth-api2-dev-forgot-password","functionVersion":"$LATEST","invokedFunctionArn":"offline_invokedFunctionArn_for_oauth-api2-dev-forgot-password","logGroupName":"offline_logGroupName_for_oauth-api2-dev-forgot-password","logStreamName":"offline_logStreamName_for_oauth-api2-dev-forgot-password","memoryLimitInMB":"128"}

	// lambdacontext.NewContext(context.Background(), &lambdacontext.LambdaContext{})

	res, err := h(context.Background(), event)
	if err != nil {
		response.Payload.Error = err.Error()
		out, _ := json.Marshal(response)
		fmt.Println(string(out))
		os.Exit(1)
	}

	response.Payload.Success = res
	out, _ := json.Marshal(response)
	fmt.Println(string(out))
}
