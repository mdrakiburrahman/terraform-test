# Storage Account example

This folder contains a [Terraform](https://www.terraform.io/) configuration that shows an example of how to 
use the [storage-account module](../../modules/data-stores/storage-account) to deploy an Azure Storage Account. 

## Pre-requisites

* Launch this `.devcontainer`
* You must have an Azure Service Principal with `Contributor` priveleges.

Please note that this code was written for Terraform 1.x+.

## Quick start

Change directory to here:
```bash
cd /workspaces/terraform-test/examples/storage-account
```

Pipe in Service Principal Creds from environment variables:

```bash
# Terraform Provider
export ARM_TENANT_ID=$spnTenantId
export ARM_CLIENT_ID=$spnClientId
export ARM_CLIENT_SECRET=$spnClientSecret
export ARM_SUBSCRIPTION_ID=$subscriptionId

# Golang Azure SDK
export AZURE_TENANT_ID=$ARM_TENANT_ID
export AZURE_CLIENT_ID=$ARM_CLIENT_ID
export AZURE_CLIENT_SECRET=$ARM_CLIENT_SECRET
export AZURE_SUBSCRIPTION_ID=$ARM_SUBSCRIPTION_ID
```

Configure module specific variables:

```bash
export TF_VAR_resource_group_name='quick-rg'
export TF_VAR_location='eastus'
export TF_VAR_tags='{ Source = "terraform", Owner = "Your Name", Project = "Messing around with terraform manually" }'
export TF_VAR_account_tier='Standard'
export TF_VAR_account_kind='StorageV2'
export TF_VAR_account_unique_postfix='8479q7h'
```

Deploy the code:

```bash
terraform init
terraform plan
terraform apply -auto-approve
```

Clean up when you're done:

```
terraform destroy -auto-approve
rm -rf .terraform
rm .terraform.lock.hcl
rm terraform.tfstate 
rm terraform.tfstate.backup
```