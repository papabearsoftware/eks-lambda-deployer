resource "aws_codebuild_project" "build" {
  name          = "example-eks-lambda-build"
  description   = "An example for deploying to EKS from CodePipeline and Lambda"
  build_timeout = "5"
  service_role  = aws_iam_role.build_role.arn

  artifacts {
    type = "CODEPIPELINE"
  }

  environment {
    compute_type                = "BUILD_GENERAL1_SMALL"
    image                       = "aws/codebuild/standard:1.0"
    type                        = "LINUX_CONTAINER"
    image_pull_credentials_type = "CODEBUILD"

    environment_variable {
      name  = "FOO"
      value = "BAR"
    }
  }

  logs_config {
    cloudwatch_logs {
      group_name  = "log-group"
      stream_name = "log-stream"
    }

    s3_logs {
      status   = "ENABLED"
      location = "${aws_s3_bucket.artifacts.id}/build-log"
    }
  }

  source {
    type            = "CODEPIPELINE"
  }
  tags = {
    Foo = "Bar"
  }
}

resource "aws_codepipeline" "pipeline" {
  name = "example-eks-lambda-pipeline"
  role_arn = aws_iam_role.codepipeline_role.arn

  artifact_store {
    location = aws_s3_bucket.artifacts.bucket
    type     = "S3"
  }

  stage {
    name = "GetSource"

    action {
      name             = "Source"
      category         = "Source"
      owner            = "ThirdParty"
      provider         = "GitHub"
      version          = "1"
      output_artifacts = ["SourceArtifact"]
      configuration = {
        Owner      = var.github_org
        Repo       = var.github_repo
        Branch     = var.github_branch
        OAuthToken = var.github_token
      }
    }
  }
    stage {
      name = "Build"

      action {
        name             = "Build"
        category         = "Build"
        owner            = "AWS"
        provider         = "CodeBuild"
        input_artifacts  = ["SourceArtifact"]
        output_artifacts = ["BuildArtifact"]
        version          = "1"

        configuration = {
          ProjectName = aws_codebuild_project.build.name
        }
      }
    }

    stage {
        name = "Deploy"

        action {
            // https://docs.aws.amazon.com/codepipeline/latest/userguide/action-reference-Lambda.html
            category = "Invoke"
            owner = "AWS"
            version = "1" // TODO make a version in labmda.tf
            name = "Lambda"
            provider = "Lambda"
            input_artifacts = ["BuildArtifact"]

            configuration = {
                FunctionName = aws_lambda_function.lambda.function_name
            }
        }
    }
}