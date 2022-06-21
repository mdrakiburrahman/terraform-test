# Automated testing

This folder contains examples of how to write automated tests for infrastructure code using Go and
[Terratest](https://terratest.gruntwork.io/).

## Pre-requisites

* Launch this `.devcontainer`
* You must have an Azure Service Principal with `Contributor` priveleges injected into this container.

## Quick start

First time build:
```bash
cd /workspaces/terraform-test/test

# Can call it whatever we want - in this case our repo name
go mod init github.com/mdrakiburrahman/terraform-test

# This creates a go.sum file with all our dependencies linked to git commits, and cleans up ones not required
go mod tidy
```

Run all the tests:

```bash
go test -v -timeout 90m
```

Run one specific test:

```bash
go test -v -timeout 90m -run 'TestStorageAccountExample'
```