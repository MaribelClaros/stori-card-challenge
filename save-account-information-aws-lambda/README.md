# save-account-information-aws-lambda

This Lambda is responsible for receiving the SNS event and creating the item in the AccountInformation table.

## Steps to run this Lambda

1. **Push Docker image**

Build the Docker image for the Lambda (using Linux AMD64 platform):

```bash

## Build image
docker build --platform linux/amd64 -t save-account-information-aws-lambda .
 
## Tag image
docker tag save-account-information-aws-lambda:latest 128624920373.dkr.ecr.us-east-1.amazonaws.com/account-information-images

## Push Image to ECR service
docker push 128624920373.dkr.ecr.us-east-1.amazonaws.com/account-information-images
```

2. Create a new Lambda function in AWS called `save-account-information-aws-lambda ` configured to use the Docker image you pushed in the previous step.