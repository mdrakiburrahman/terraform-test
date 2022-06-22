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

Ideally, the workflow would look more like this - let's say we had an `App` that depended on a back-end `MySQL` database.

We want to iterate on the `App` but don't want to touch the `MySQL` since it takes long time to deploy.

1. Run terraform apply on the `mysql` module.
2. Run terraform apply on the `hello-world-app` module.
3. Now, you start doing iterative development:
    a. Make a change to the `hello-world-app` module - which shares State with `mysql` - see [here](https://github.com/brikis98/terraform-up-and-running-code/blob/503ab1f5055917f2d0c715a6b1aa0b9dfb716354/code/terraform/09-testing-terraform-code/test/hello_world_integration_test.go#L73)
    b. Rerun terraform apply on the `hello-world-app` module to deploy your updates - gets persisted in the same state file as before
    c. Run validations to make sure everything is working.
    d. If everything works, move on to the next step. If not, go back to step (3a).
4. Run terraform destroy on the `hello-world-app` module.
5. Run terraform destroy on the `mysql` module

We want to speed up 3a - 3d to make our lives easier.

| **Note:** This only works because we are not changing our backend state file over and over. If we did, Terraform would break and complain. So for example in our case, we set the `uniqueID` in the statefile, so throughout the development lifecycle we need to make sure it doesn't change. When we're done our work for the day we can redeploy again next, day, which will recreate the infra alongside `uniqueID`. We can then choose to delete the stale state files locally and remotely I guess.

Example:
```bash
# Blow away old local state
rm -rf /workspaces/terraform-test/examples/storage-account/.terraform
rm -rf /workspaces/terraform-test/examples/storage-account/.test-data
rm -rf /workspaces/terraform-test/examples/storage-account/.terraform.lock.hcl

# 1. Deploy one-time
SKIP_teardown_storageAccount=true \
SKIP_validate_storageAccount=true \
go test -timeout 30m -run 'TestStorageAccountExampleWithStages'
# PASS
# ok      github.com/mdrakiburrahman/terraform-test       93.488s

# 3. a, b, c - Validate
# `SKIP_deploy_storageAccount` - this must be set! Otherwise the new uniqueID generator will force generate a new Statefile in Blob, which will confuse Terraform versus the Terratest local copy - which references the State file
# I guess we could just store state locally to avoid this - but I think that'll still redeploy the Storage Account because of the uniqueID. So basically, it's not worth it.
SKIP_deploy_storageAccount=true \
SKIP_teardown_storageAccount=true \
go test -timeout 30m -run 'TestStorageAccountExampleWithStages'
# PASS
# ok      github.com/mdrakiburrahman/terraform-test       8.912s

# Note that if we make one of our tests fail on purpose, only then do we see the test breakdown:
# --- FAIL: TestStorageAccountExampleWithStages (9.01s)
#     --- FAIL: TestStorageAccountExampleWithStages/storage_account_location_tf_input_matched_arm_output (0.00s)
#         storage_account_example_test.go:129: 
#                 Error Trace:    storage_account_example_test.go:129
#                 Error:          Should not be: "canadacentral"
#                 Test:           TestStorageAccountExampleWithStages/storage_account_location_tf_input_matched_arm_output
#                 Messages:       Storage Account Location TF Input = ARM Output
# FAIL
# exit status 1
# FAIL    github.com/mdrakiburrahman/terraform-test       9.036s



# The beauty is:
# 1. We can rerun validations super quick
# 2. If we had another module that dependend on the Storage Account, that would run quickly too!

# 5 - Destroy
SKIP_deploy_storageAccount=true \
SKIP_validate_storageAccount=true \
go test -timeout 30m -run 'TestStorageAccountExampleWithStages'
# PASS
# ok      github.com/mdrakiburrahman/terraform-test       144.295s
```