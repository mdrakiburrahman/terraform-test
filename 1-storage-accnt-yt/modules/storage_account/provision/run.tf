terraform {
  required_version = "~> 1.0"
  required_providers {
    azurerm = "~> 3.9.0"
  }
}

provider "azurerm" {
  features {}
}

module "core_storage" {
  source              = "../"
  resource_group_name = var.resource_group_name
  location            = var.location
  account_tier        = var.account_tier
  unique_postfix      = var.unique_postfix
  account_kind        = var.account_kind
}

output "storage_output" {
  value = {
    storage_account_name = module.core_storage.storage.name
    resource_group_name  = module.core_storage.resource_group_name
    account_tier         = module.core_storage.account_tier
    account_kind         = module.core_storage.account_kind
  }
}

output "storage_account_name" {
  value = module.core_storage.storage.name
}

output "resource_group_name" {
  value = module.core_storage.storage.resource_group_name
}
