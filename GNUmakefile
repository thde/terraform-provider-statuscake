WEBSITE_REPO=github.com/hashicorp/terraform-website
PKG_NAME=statuscake

default: build

lint:
	golangci-lint run

build: lint
	go install

test:
	go test ./...

website-test:
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs validate

.PHONY: build lint test website-test

