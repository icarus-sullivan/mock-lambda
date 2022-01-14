package lambda

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

type ApiHandler func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

type Response struct {
	Payload Payload `json:"offline_payload"`
}

type Payload struct {
	Error   string                         `json:"error,omitempty"`
	Success events.APIGatewayProxyResponse `json:"success,omitempty"`
}

func sanitize(json string) string {
	// When we get this env value it's escaped in a way javascript understands but not
	// go
	// 	- unescape to a single backslash
	return strings.ReplaceAll(json, "\\\\", "\\")
}

func Start(h ApiHandler) {
	response := Response{
		Payload: Payload{},
	}

	sanitizedEvent := sanitize(os.Getenv("LAMBDA_EVENT"))

	var event events.APIGatewayProxyRequest
	err := json.Unmarshal([]byte(sanitizedEvent), &event)
	if err != nil {
		response.Payload.Error = err.Error()
		out, _ := json.Marshal(response)
		fmt.Println(string(out))
		os.Exit(1)
	}

	contextEvent := sanitize(os.Getenv("LAMBDA_CONTEXT"))
	ctxArgs := map[string]interface{}{}
	err = json.Unmarshal([]byte(contextEvent), &ctxArgs)
	if err != nil {
		response.Payload.Error = err.Error()
		out, _ := json.Marshal(response)
		fmt.Println(string(out))
		os.Exit(1)
	}

	ctx := lambdacontext.NewContext(context.Background(), &lambdacontext.LambdaContext{
		AwsRequestID:       ctxArgs["awsRequestId"].(string),
		InvokedFunctionArn: ctxArgs["invokedFunctionArn"].(string),
	})

	res, err := h(ctx, event)
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
