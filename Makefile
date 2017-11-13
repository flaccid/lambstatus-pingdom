DOCKER_REGISTRY = index.docker.io
IMAGE_NAME = checks2metrics
IMAGE_VERSION = latest
IMAGE_ORG = flaccid
IMAGE_TAG = $(DOCKER_REGISTRY)/$(IMAGE_ORG)/$(IMAGE_NAME):$(IMAGE_VERSION)

WORKING_DIR := $(shell pwd)
ROLE_POLICY := $(shell cat role_policy.json)

.DEFAULT_GOAL := help

.PHONY: docker-build docker-push

build:: ## Builds the checks2metrics binary
		@go build -o bin/checks2metrics cli/checks2metrics.go

build-static-linux:: ## Builds a static linux binary
		@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
			go build \
			-o bin/checks2metrics \
			-a -ldflags '-extldflags "-static"' \
				cli/checks2metrics.go

docker-build:: ## Builds the docker image locally
		@docker build --pull \
		-t $(IMAGE_TAG) $(WORKING_DIR)

docker-push:: ## Pushes the docker image to the registry
		@docker push $(IMAGE_TAG)

docker-release:: docker-build docker-push ## Builds and pushes the docker image to the registry

lambda-attach-iam-role-policy:: ## Attaches the IAM role policy
		aws iam attach-role-policy \
  		--role-name lambda_basic_execution \
  		--policy-arn arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole \
  			|| exit 1

lambda-create-iam-role:: ## Creates the AWS IAM role for Lambda
		aws iam create-role \
  		--role-name lambda_basic_execution \
  		--assume-role-policy-document '$(ROLE_POLICY)' || exit 1

lambda-create-function:: ## Creates the AWS Lambda function
		aws lambda create-function \
  		--function-name checks2metrics \
  		--zip-file fileb://handler.zip \
  		--role $(shell aws iam get-role --role-name lambda_basic_execution --query 'Role.Arn' --output text) \
  		--runtime python2.7 \
  		--handler handler.Handle || exit 1

lambda-invoke-function:: ## Invokes the AWS Lambda function
		aws lambda invoke \
  		--function-name checks2metrics \
  		--invocation-type RequestResponse \
  		--log-type Tail  /dev/stderr \
  		--query 'LogResult' \
  		--output text

run:: ## Runs the executable
		bin/checks2metrics

# A help target including self-documenting targets (see the awk statement)
define HELP_TEXT
Usage: make [TARGET]... [MAKEVAR1=SOMETHING]...

Available targets:
endef
export HELP_TEXT
help: ## This help target
	@cat .banner
	@echo
	@echo "$$HELP_TEXT"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / \
		{printf "\033[36m%-30s\033[0m  %s\n", $$1, $$2}' $(MAKEFILE_LIST)
