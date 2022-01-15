package lambda

import (
	"fmt"
	"os"
)

func Start(h interface{}) {
	fmt.Println("IS_LAMBDA_AUTHORIZER", os.Getenv("IS_LAMBDA_AUTHORIZER"))
	// Most common - probably
	if os.Getenv("IS_LAMBDA_AUTHORIZER") != "true" {
		ApiHandler(h.(ApiLambda))
	}
	if os.Getenv("IS_LAMBDA_AUTHORIZER") == "true" {
		AuthorizerHandler(h.(AuthorizerLambda))
	}
}
