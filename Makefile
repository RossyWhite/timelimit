include .env

GOCMD := go
AWSCMD := aws
GOBUILD :=$(GOCMD) build
GOTEST :=$(GOCMD) test
STACK_NAME := $(STACK_NAME)
BUILDDIR := ./build
SRCDIR := ./src
TEMPLATE_FILE := 'template.yml'
BUCKET_NAME := $(STACK_NAME)
BORN := $(BORN)
LIFETIME := $(LIFETIME)
URL := $(URL)
INTERVAL := $(INTERVAL)
TZ := $(TZ)

.PHONY: build
build: clean
	@GOARCH=amd64 GOOS=linux GO111MODULE=on $(GOBUILD) -o $(BUILDDIR)/$(STACK_NAME) $(SRCDIR)/main.go

.PHONY: clean
clean:
	@rm -rf $(BUILDDIR)/*

.PHONY: bucket
bucket:
	@if ! $(AWSCMD) s3api head-bucket --bucket $(BUCKET_NAME) 2>/dev/null ; then $(AWSCMD) s3 mb s3://$(BUCKET_NAME); fi

.PHONY: package
package: clean build bucket
	@$(AWSCMD) cloudformation package \
    --template-file $(TEMPLATE_FILE) \
    --s3-bucket $(BUCKET_NAME) \
    --output-template-file $(BUILDDIR)/$(TEMPLATE_FILE)

.PHONY: deploy
deploy: package
	@$(AWSCMD) cloudformation deploy \
    --template-file $(BUILDDIR)/$(TEMPLATE_FILE) \
    --stack-name $(STACK_NAME) \
    --capabilities CAPABILITY_IAM \
    --parameter-overrides \
    	"born"="$(BORN)" \
    	"lifetime"="$(LIFETIME)" \
    	"url"="$(URL)" \
    	"interval"="$(INTERVAL)" \
    	"tz"="$(TZ)"

.PHONY: delete
delete:
	@$(AWSCMD) cloudformation delete-stack \
	--stack-name $(STACK_NAME)

	@$(AWSCMD) s3 rm s3://$(BUCKET_NAME) --recursive
	@$(AWSCMD) s3 rb s3://$(BUCKET_NAME)

.PHONY: test
test:
	@export BORN=$(BORN) LIFETIME=$(LIFETIME) URL=$(URL); $(GOTEST) $(SRCDIR)
