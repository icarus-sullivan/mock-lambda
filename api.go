package lambda

import (
	"os"
)

func Start(h interface{}) {

	// Most common - probably
	if os.Getenv("IS_LAMBDA_AUTHORIZER") != "true" {
		ApiHandler(h.(ApiLambda))
	}
	if os.Getenv("IS_LAMBDA_AUTHORIZER") == "true" {
		AuthorizerHandler(h.(AuthorizerLambda))
	}
}
