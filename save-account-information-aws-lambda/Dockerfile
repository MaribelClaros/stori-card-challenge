# Stage 1: Build the Go binary on amd64 platform
FROM  golang:1.18 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

# Build de la app
RUN go build -o main ./cmd

FROM  public.ecr.aws/lambda/go:1

COPY --from=builder /app/main ${LAMBDA_TASK_ROOT}/main
COPY --from=builder /app/aws_config.json /var/task/


CMD ["main"]
