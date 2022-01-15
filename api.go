package lambda

import (
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

func Api(h func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) Response {
	response := Response{
		Payload: Payload{},
	}

	sanitizedEvent := SanitizeJSON(os.Getenv("LAMBDA_EVENT"))
	event := events.APIGatewayProxyRequest{}
	err := json.Unmarshal([]byte(sanitizedEvent), &event)
	if err != nil {
		response.Payload.Error = err.Error()
		return response
	}

	contextEvent := SanitizeJSON(os.Getenv("LAMBDA_CONTEXT"))
	ctxArgs := map[string]interface{}{}
	err = json.Unmarshal([]byte(contextEvent), &ctxArgs)
	if err != nil {
		response.Payload.Error = err.Error()
		return response
	}

	ctx := lambdacontext.NewContext(context.Background(), &lambdacontext.LambdaContext{
		AwsRequestID:       ctxArgs["awsRequestId"].(string),
		InvokedFunctionArn: ctxArgs["invokedFunctionArn"].(string),
	})

	res, err := h(ctx, event)
	if err != nil {
		response.Payload.Error = err.Error()
		return response
	}

	response.Payload.Success = res
	return response
}
