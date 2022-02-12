package lambda

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"testing"
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

	err := ml.start(func(ctx context.Context, request events.SQSEvent) error {
		return nil
	})

	if err != nil {
		t.Error(err)
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

	err := ml.start(func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return events.APIGatewayProxyResponse{}, nil
	})

	if err != nil {
		t.Error(err)
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

	err := ml.start(func(request events.APIGatewayCustomAuthorizerRequestTypeRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
		return events.APIGatewayCustomAuthorizerResponse{}, nil
	})

	if err != nil {
		t.Error(err)
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

	err := ml.start(func(request events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
		return events.APIGatewayCustomAuthorizerResponse{}, nil
	})

	if err != nil {
		t.Error(err)
		t.Fail()
	}
}
