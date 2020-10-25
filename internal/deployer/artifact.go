package deployer

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/papabearsoftware/eks-lambda-deployer/internal/util"
)

type DeploymentJSON struct {
	Cluster        string `json:"cluster"`
	RollbackOnFail bool   `json:"rollback_on_fail"`
	Deployment     string `json:"deployment"`
	Tag            string `json:"tag"`
	Namespace      string `json:"namespace"`
}

func retrieveS3Artifact(bucket string, key string) error {

	f, err := os.Create(fmt.Sprintf("/tmp/artifact-%s.zip", jobID))

	if err != nil {
		util.LogError(fmt.Sprintf("Error creating '/tmp/artifact-%s.zip'", jobID), err.Error())
		return err
	}

	defer f.Close()

	downloader := s3manager.NewDownloader(sess)

	util.LogInfo(fmt.Sprintf("Retrieving input artifact from %s/%s", bucket, key))

	_, err = downloader.Download(f,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})

	if err != nil {
		util.LogError(fmt.Sprintf("Error downloading %s/%s", bucket, key), err.Error())
		return err
	}

	util.LogDebug("Downloaded input artifact")

	return nil
}

func parseArtifact() error {
	util.LogDebug("Starting artifact parsing")

	r, err := zip.OpenReader(fmt.Sprintf("/tmp/artifact-%s.zip", jobID))

	if err != nil {
		util.LogError(fmt.Sprintf("Error opening '/tmp/artifact-%s.zip'", jobID), err.Error())
		return err
	}

	defer r.Close()

	util.LogDebug("Created zip reader")

	configExists := false

	for _, f := range r.File {
		if f.Name == "deployer_config.json" {
			util.LogDebug("Iterating through files in archive")

			configExists = true

			rc, err := f.Open()

			util.LogDebug("Opened deployer_config.json")

			if err != nil {
				util.LogError("Error opening config file", err.Error())
				return err
			}

			body, err := ioutil.ReadAll(rc)

			util.LogDebug("Read file body into []byte")

			if err != nil {
				util.LogError("Error reading config file", err.Error())
				return err
			}

			err = json.Unmarshal(body, &deploymentJSON)

			if err != nil {
				util.LogError("Error umarshaling config file contents into struct", err.Error())
				return err
			}

			util.LogDebug("Unmarshalled successfully")
			rc.Close()
		}
	}

	if !configExists {
		util.LogError("Config file not found!", "")
		return errors.New("Config file not found")
	}

	util.LogDebug("Returning nil error")

	return nil
}
