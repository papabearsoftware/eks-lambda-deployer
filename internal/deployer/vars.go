package deployer

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws/session"
)

var (
	sess           *session.Session
	err            error
	jobID          string
	deploymentJSON DeploymentJSON

	// Errors
	//ErrConfigNotFound is thrown when the configuration file is not found
	ErrConfigNotFound = errors.New("Config file not found")
)
