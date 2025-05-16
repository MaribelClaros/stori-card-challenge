# process-transactions-aws-lambda

This Lambda function is responsible for processing transactions, sending emails, and invoking SNS.

## Steps to run this Lambda

1. **Push Docker image**

Build the Docker image for the Lambda (using Linux ARM64 platform):

```bash

## Build image
 docker build --platform linux/amd64 -t process-transactions-aws-lambda . 
 
## Tag image
docker tag process-transactions-aws-lambda:latest 128624920373.dkr.ecr.us-east-1.amazonaws.com/docker-lambda-images

## Push Image to ECR service
docker push 128624920373.dkr.ecr.us-east-1.amazonaws.com/docker-lambda-images
```

2. Create a new Lambda function in AWS called `process-transactions-aws-lambda ` configured to use the Docker image you pushed in the previous step.