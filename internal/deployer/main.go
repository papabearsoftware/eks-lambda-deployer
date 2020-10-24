package deployer

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/papabearsoftware/eks-lambda-deployer/internal/util"
)

func LambdaHandler(ctx context.Context, event events.CodePipelineEvent) {
	jobID = event.CodePipelineJob.ID
	util.LogInfo(fmt.Sprintf("Received deployment ID: %s", jobID))

	sess, err = session.NewSession(&aws.Config{
		Region:     aws.String(os.Getenv("AWS_REGION")),
		MaxRetries: aws.Int(3)},
	)

	if err != nil {
		util.LogError("Error creating AWS Session", err.Error())
		markDeploymentFailure("Credentials Error", "Error creating AWS Session")
	}

	err = retrieveS3Artifact(event.CodePipelineJob.Data.InputArtifacts[0].Location.S3Location.BucketName,
		event.CodePipelineJob.Data.InputArtifacts[0].Location.S3Location.ObjectKey)

	if err != nil {
		markDeploymentFailure("Artifact Retrieval Error", "Error retrieving input artifact")
	}

	err = parseArtifact()

	if err != nil {
		markDeploymentFailure("Config Parse Error", "Error parsing config file from input artifact")
	}

	// Deploy was a success!
	util.LogInfo("Marking job successful")
	markDeploymentSuccess(event.CodePipelineJob.Data.ContinuationToken)
}
