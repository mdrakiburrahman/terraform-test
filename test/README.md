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

Run all the test modules:

```bash
go test -v -timeout 90m
```

Run one specific modules and all test cases within it:

```bash
# Blow away old local state
rm -rf /workspaces/terraform-test/examples/storage-account/.terraform
rm -rf /workspaces/terraform-test/examples/storage-account/.test-data
rm -rf /workspaces/terraform-test/examples/storage-account/.terraform.lock.hcl

go test -v -timeout 90m -run 'TestStorageAccountExampleWithStages'
```

Sample output as follows for the 3 stages:
```bash
#
# Deploy
#
# TestStorageAccountExampleWithStages 2022-06-21T22:19:29Z test_structure.go:27: The 'SKIP_deploy_storageAccount' environment variable is not set, so executing stage 'deploy_storageAccount'.
# TestStorageAccountExampleWithStages 2022-06-21T22:19:29Z save_test_data.go:188: Storing test data in ../examples/storage-account/.test-data/TerraformOptions.json so it can be reused later
# ...
# TestStorageAccountExampleWithStages 2022-06-21T22:21:04Z logger.go:66: module.storage_account.azurerm_storage_account.storage_accnt: Creation complete after 26s [id=/subscriptions/ce859648-30e1-4135-9d0f-8358aebfe789/resourceGroups/terratest-storage-account-yyxcun/providers/Microsoft.Storage/storageAccounts/storageacctyyxcun]
# TestStorageAccountExampleWithStages 2022-06-21T22:21:04Z logger.go:66: 
# TestStorageAccountExampleWithStages 2022-06-21T22:21:04Z logger.go:66: Apply complete! Resources: 2 added, 0 changed, 0 destroyed.
#
# Test
#
# ...
# TestStorageAccountExampleWithStages 2022-06-21T22:21:04Z test_structure.go:27: The 'SKIP_validate_storageAccount' environment variable is not set, so executing stage 'validate_storageAccount'.
# TestStorageAccountExampleWithStages 2022-06-21T22:21:04Z save_test_data.go:215: Loading test data from ../examples/storage-account/.test-data/TerraformOptions.json
# === RUN   TestStorageAccountExampleWithStages/storage_account_location_tf_input_matched_tf_output
# === RUN   TestStorageAccountExampleWithStages/storage_account_location_tf_input_matched_arm_output
#
# Destroy
#
# ...
# TestStorageAccountExampleWithStages 2022-06-21T22:21:14Z test_structure.go:27: The 'SKIP_teardown_storageAccount' environment variable is not set, so executing stage 'teardown_storageAccount'.
# TestStorageAccountExampleWithStages 2022-06-21T22:21:14Z save_test_data.go:215: Loading test data from ../examples/storage-account/.test-data/TerraformOptions.json
# ...
# TestStorageAccountExampleWithStages 2022-06-21T22:22:35Z logger.go:66: Destroy complete! Resources: 2 destroyed.
# TestStorageAccountExampleWithStages 2022-06-21T22:22:35Z logger.go:66: 
# --- PASS: TestStorageAccountExampleWithStages (186.32s)
#     --- PASS: TestStorageAccountExampleWithStages/storage_account_location_tf_input_matched_tf_output (0.00s)
#     --- PASS: TestStorageAccountExampleWithStages/storage_account_location_tf_input_matched_arm_output (0.00s)
# PASS
# ok      github.com/mdrakiburrahman/terraform-test       186.357s
```

So this sample run takes about 3.5 mins.

## Development workflow via `Stages`

The actual development workflow can be fasterm since the code supports stage tags.

Ideally, the workflow would look more like this:
1. Run `terraform apply`.
2. Now, you start doing iterative development:
    a. Make a change to the example/root modules' `.tf` files as needed.
    b. Rerun terraform apply on the to deploy your updates.
    c. Run validations to make sure everything is working.
    d. If everything works, move on to the next step. If not, go back to step (2a).
3. Run `terraform destroy` on both modules.

We want to speed up 2a - 2d to make our lives easier.

```bash
# 1.
SKIP_teardown_storageAccount=true \
SKIP_validate_storageAccount=true \
go test -timeout 30m -run 'TestStorageAccountExampleWithStages'
# PASS
# ok      github.com/mdrakiburrahman/terraform-test       93.488s

# 2. b, c
SKIP_deploy_storageAccount=true \
SKIP_validate_storageAccount=true \
go test -timeout 30m -run 'TestStorageAccountExampleWithStages'

# 2 c only

```