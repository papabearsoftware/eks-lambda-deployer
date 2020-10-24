package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/papabearsoftware/eks-lambda-deployer/internal/deployer"
	"github.com/papabearsoftware/eks-lambda-deployer/internal/util"
)

// main only serves to start the lambda
func main() {
	util.LogInfo("Deployer booted")
	lambda.Start(deployer.LambdaHandler)
}
