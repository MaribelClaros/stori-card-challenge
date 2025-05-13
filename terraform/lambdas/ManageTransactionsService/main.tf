# Create IAM role for the Lambda function
resource "aws_iam_role" "lambda_execution_role" {
  name = "lambda_csv_processor_role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [{
      Action    = "sts:AssumeRole",
      Effect    = "Allow",
      Principal = {
        Service = "lambda.amazonaws.com"
      }
    }]
  })
}

# Attach basic execution and log permissions to the role
resource "aws_iam_role_policy_attachment" "lambda_basic_execution" {
  role       = aws_iam_role.lambda_execution_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

# (Optional) Attach permissions to access DynamoDB, S3, or SNS
# Add additional inline or managed policies here if needed

# Create the Lambda function from ECR image
resource "aws_lambda_function" "csv_processor" {
  function_name = "csv-processor"
  role          = aws_iam_role.lambda_execution_role.arn
  package_type  = "Image"
  image_uri     = "${aws_ecr_repository.csv_processor.repository_url}:latest" # Replace tag if needed
  timeout       = 30
  memory_size   = 512

  environment {
    variables = {
      DYNAMODB_TABLE_NAME = aws_dynamodb_table.transactions.name
      SNS_TOPIC_ARN       = aws_sns_topic.notifications.arn
    }
  }
}
