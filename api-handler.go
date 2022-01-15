package lambda

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

type ApiResponse struct {
	ApiPayload ApiPayload `json:"offline_payload"`
}

type ApiPayload struct {
	Error   string                         `json:"error,omitempty"`
	Success events.APIGatewayProxyResponse `json:"success,omitempty"`
}

func ApiHandler(h func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) {
	response := ApiResponse{
		ApiPayload: ApiPayload{},
	}

	sanitizedEvent := SanitizeJSON(os.Getenv("LAMBDA_EVENT"))

	event := events.APIGatewayProxyRequest{}
	err := json.Unmarshal([]byte(sanitizedEvent), &event)
	if err != nil {
		response.ApiPayload.Error = err.Error()
		out, _ := json.Marshal(response)
		fmt.Println(string(out))
		os.Exit(1)
	}

	contextEvent := SanitizeJSON(os.Getenv("LAMBDA_CONTEXT"))
	ctxArgs := map[string]interface{}{}
	err = json.Unmarshal([]byte(contextEvent), &ctxArgs)
	if err != nil {
		response.ApiPayload.Error = err.Error()
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
		response.ApiPayload.Error = err.Error()
		out, _ := json.Marshal(response)
		fmt.Println(string(out))
		os.Exit(1)
	}

	response.ApiPayload.Success = res
	out, _ := json.Marshal(response)
	fmt.Println(string(out))
}
