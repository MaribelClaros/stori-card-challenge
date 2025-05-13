# Create a unique S3 bucket for uploading CSV files
resource "aws_s3_bucket" "csv_upload_bucket" {
  bucket = "transactions-bucket"

  tags = {
    Name        = "CSV Transactions Bucket"
    Environment = "dev"
  }
}

# Block all public access to the bucket
resource "aws_s3_bucket_public_access_block" "block_public_access" {
  bucket = aws_s3_bucket.csv_upload_bucket.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

# Allow S3 to invoke the Lambda when a new CSV file is uploaded
resource "aws_lambda_permission" "allow_s3_to_invoke_lambda" {
  statement_id  = "AllowS3Invoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.csv_processor.function_name
  principal     = "s3.amazonaws.com"
  source_arn    = aws_s3_bucket.csv_upload_bucket.arn
}

# Configure the bucket to trigger the Lambda on CSV file upload
resource "aws_s3_bucket_notification" "s3_to_lambda_notification" {
  bucket = aws_s3_bucket.csv_upload_bucket.id

  lambda_function {
    lambda_function_arn = aws_lambda_function.csv_processor.arn
    events              = ["s3:ObjectCreated:*"]
    filter_prefix       = "incoming/"
    filter_suffix       = ".csv"
  }

  # Wait until permission is set before applying this resource
  depends_on = [aws_lambda_permission.allow_s3_to_invoke_lambda]
}
