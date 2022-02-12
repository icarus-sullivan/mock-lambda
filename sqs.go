package lambda

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

func sqs(h func(ctx context.Context, request events.SQSEvent) error) Response {
	response := Response{
		Payload: Payload{},
	}

	event := events.SQSEvent{}
	err := decode(os.Getenv("LAMBDA_EVENT"), &event)
	if err != nil {
		response.Payload.Error = err.Error()
		return response
	}

	ctxArgs := map[string]interface{}{}
	err = decode(os.Getenv("LAMBDA_CONTEXT"), &ctxArgs)
	if err != nil {
		response.Payload.Error = err.Error()
		return response
	}

	ctx := lambdacontext.NewContext(context.Background(), &lambdacontext.LambdaContext{
		AwsRequestID:       ctxArgs["awsRequestId"].(string),
		InvokedFunctionArn: ctxArgs["invokedFunctionArn"].(string),
	})

	err = h(ctx, event)
	if err != nil {
		response.Payload.Error = err.Error()
		return response
	}

	response.Payload.Success = err != nil
	return response
}
