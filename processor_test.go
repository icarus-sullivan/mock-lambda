package lambda

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestStart_DeterminesSQSInvocation(t *testing.T) {
	ml := MockLambda{api: func(h func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) Response {
		t.Errorf("incorrect invocation")
		t.Fail()
		return Response{}
	}, request: func(h func(request events.APIGatewayCustomAuthorizerRequestTypeRequest) (events.APIGatewayCustomAuthorizerResponse, error)) Response {
		t.Errorf("incorrect invocation")
		t.Fail()
		return Response{}
	}, token: func(h func(request events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error)) Response {
		t.Errorf("incorrect invocation")
		t.Fail()
		return Response{}
	}, sqs: func(h func(ctx context.Context, request events.SQSEvent) error) Response {
		return Response{}
	}}

	response := ml.start(func(ctx context.Context, request events.SQSEvent) error {
		return nil
	})

	if response.Payload.Error != "" {
		t.Error(response.Payload.Error)
		t.Fail()
	}
}

func TestStart_DeterminesAPIInvocation(t *testing.T) {
	ml := MockLambda{api: func(h func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) Response {
		return Response{}
	}, request: func(h func(request events.APIGatewayCustomAuthorizerRequestTypeRequest) (events.APIGatewayCustomAuthorizerResponse, error)) Response {
		t.Errorf("incorrect invocation")
		t.Fail()
		return Response{}
	}, token: func(h func(request events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error)) Response {
		t.Errorf("incorrect invocation")
		t.Fail()
		return Response{}
	}, sqs: func(h func(ctx context.Context, request events.SQSEvent) error) Response {
		t.Errorf("incorrect invocation")
		t.Fail()
		return Response{}
	}}

	response := ml.start(func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return events.APIGatewayProxyResponse{}, nil
	})

	if response.Payload.Error != "" {
		t.Error(response.Payload.Error)
		t.Fail()
	}
}

func TestStart_DeterminesRequestInvocation(t *testing.T) {
	ml := MockLambda{api: func(h func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) Response {
		t.Errorf("incorrect invocation")
		t.Fail()
		return Response{}
	}, request: func(h func(request events.APIGatewayCustomAuthorizerRequestTypeRequest) (events.APIGatewayCustomAuthorizerResponse, error)) Response {
		return Response{}
	}, token: func(h func(request events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error)) Response {
		t.Errorf("incorrect invocation")
		t.Fail()
		return Response{}
	}, sqs: func(h func(ctx context.Context, request events.SQSEvent) error) Response {
		t.Errorf("incorrect invocation")
		t.Fail()
		return Response{}
	}}

	response := ml.start(func(request events.APIGatewayCustomAuthorizerRequestTypeRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
		return events.APIGatewayCustomAuthorizerResponse{}, nil
	})

	if response.Payload.Error != "" {
		t.Error(response.Payload.Error)
		t.Fail()
	}
}

func TestStart_DeterminesTokenInvocation(t *testing.T) {
	ml := MockLambda{api: func(h func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) Response {
		t.Errorf("incorrect invocation")
		t.Fail()
		return Response{}
	}, request: func(h func(request events.APIGatewayCustomAuthorizerRequestTypeRequest) (events.APIGatewayCustomAuthorizerResponse, error)) Response {
		t.Errorf("incorrect invocation")
		t.Fail()
		return Response{}
	}, token: func(h func(request events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error)) Response {
		return Response{}
	}, sqs: func(h func(ctx context.Context, request events.SQSEvent) error) Response {
		t.Errorf("incorrect invocation")
		t.Fail()
		return Response{}
	}}

	response := ml.start(func(request events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
		return events.APIGatewayCustomAuthorizerResponse{}, nil
	})

	if response.Payload.Error != "" {
		t.Error(response.Payload.Error)
		t.Fail()
	}
}
