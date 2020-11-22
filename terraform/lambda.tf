
resource "aws_lambda_function" "lambda" {
  filename      = "../latest.zip" // TODO get from GitHub
  function_name = "eks-deployer"
  role          = aws_iam_role.lambda.arn
  handler       = "main" // binary is named main
  memory_size = 256
  timeout = 600 // 5 minutes
  source_code_hash = filebase64sha256("../latest.zip") // TODO change to correct local file path

  runtime = "go1.x"

  environment {
    variables = {
      foo = "bar"
    }
  }
}