# eks-lambda-deployer

A Lambda function to deploy updates to deployments in EKS as a CodePipeline step.

**THIS IS A WORK IN PROGRESS**

## How It Works

The Lambda function will read an incoming CodePipeline event to determine where the artifact containing the new container tag.

The pipeline flow is:

- CodePipeline retrieves latest code

- CodeBuild builds the new container image and writes a simple JSON to a file that gets uploaded to S3. This file is the build artifact that the Lambda uses

- The Lambda pulls down the artifact, parses it, and deploys it. In the case of a failure, it rolls back if `rollback_on_fail` is `true`

**Note:** You can set the `DEBUG` environment variable to "true" for most detailed logging


## JSON

**Note:** For now, the file name must be `deployer_config.json`

```json
{
    "cluster": "eks-staging",
    "rollback_on_fail": true,
    "deployment": "sample-app",
    "namespace": "default",
    "containers": [
        {
            "tag": "1603656257",
            "container_name": "app"
        }
    ]    
}
```

`cluster` must be a context defined in the kubernetes config file

## TODO

Add ability to add containers to running pod

Write sample Terraform