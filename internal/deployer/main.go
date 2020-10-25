package deployer

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/papabearsoftware/eks-lambda-deployer/internal/util"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
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

	// We've parsed the deployment json, now time to deploy

	// Base64 encoded kubeconfig
	encodedKubeConfig := os.Getenv("KUBECONFIG")

	if encodedKubeConfig == "" {
		util.LogError("KUBECONFIG env var not set", "No Kubeconfig env var")
		markDeploymentFailure("Missing Kubeconfig", "KUBECONFIG environment variable not set")
	}

	// Decode kubeconfig from base64 string to []byte
	kubeConfigBytes, err := base64.StdEncoding.DecodeString(encodedKubeConfig)

	if err != nil {
		util.LogError("Error decoding kubeconfig from base64", err.Error())
		markDeploymentFailure("Decoding Error", "Error decoding kubeconfig from base64")
	}

	f, err := os.Create("/tmp/kubeconfig")

	if err != nil {
		util.LogError("Error creating /tmp/kubeconifg", err.Error())
		markDeploymentFailure("Error creating file", "Error creating /tmp/kubeconfig")
	}

	// Write bytes to file
	_, err = f.Write(kubeConfigBytes)

	if err != nil {
		f.Close()
		util.LogError("Error writing /tmp/kubeconifg", err.Error())
		markDeploymentFailure("Error writing file", "Error writing /tmp/kubeconfig")
	}

	f.Close()

	// Create kubernetes client config using the "cluster" from the deploymentJSON as the context
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: "/tmp/kubeconfig"}, &clientcmd.ConfigOverrides{
			CurrentContext: deploymentJSON.Cluster,
		}).ClientConfig()

	if err != nil {
		util.LogError("Error creating kubernetes config from /tmp/kubeconfig", err.Error())
		markDeploymentFailure("Error creating kubernetes config", "Error creating kubernetes config")
	}

	cs, err := kubernetes.NewForConfig(config)

	if err != nil {
		util.LogError("Error creating clientset", err.Error())
		markDeploymentFailure("Error creating clientset", "Error creating clientset")
	}

	kube = KubeClient{
		Client: cs,
	}

	err = deploy()

	if err != nil {
		util.LogError("Returned error from deploy()", err.Error())
	}

	// Deploy was a success!
	util.LogInfo("Marking job successful")
	markDeploymentSuccess(event.CodePipelineJob.Data.ContinuationToken)
}
