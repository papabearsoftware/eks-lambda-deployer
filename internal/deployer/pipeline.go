package deployer

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/codepipeline"
	"github.com/papabearsoftware/eks-lambda-deployer/internal/util"
)

func markDeploymentSuccess(continuationToken string) {
	cp := codepipeline.New(sess)

	i := &codepipeline.PutJobSuccessResultInput{
		JobId: aws.String(jobID),
	}

	if continuationToken != "" {
		i.ContinuationToken = aws.String(continuationToken)
	}

	o, err := cp.PutJobSuccessResult(i)

	if err != nil {
		util.LogError(fmt.Sprintf("Error putting job success result. Response: %v", o), err.Error())
	} else {
		util.LogInfo(fmt.Sprintf("PutJobSuccessResult was successful. Response: %v", o))
	}

}

func markDeploymentFailure(failureType string, failureMessage string) {
	cp := codepipeline.New(sess)

	d := &codepipeline.FailureDetails{
		Message: aws.String(failureMessage),
		Type:    aws.String(failureType),
	}

	o, err := cp.PutJobFailureResult(&codepipeline.PutJobFailureResultInput{
		JobId:          aws.String(jobID),
		FailureDetails: d,
	})

	if err != nil {
		util.LogError(fmt.Sprintf("Error putting job failure result. Response: %v", o), err.Error())
	} else {
		util.LogInfo(fmt.Sprintf("PutJobFailureResult was successful. Response: %v", o))
	}

}
