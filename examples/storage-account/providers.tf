terraform {
  
  required_version = "~> 1.0"
  required_providers {
    azurerm = "~> 3.9.0"
  }

  backend "azurerm" {
    # This backend configuration is filled in automatically at test time by Terratest. 
    # If you wish to run this example manually - see README for how to pass in env variables.

    # storage_account_name = "abcd1234"
    # container_name       = "tfstate"
    # key                  = "prod.terraform.tfstate"
    # access_key           = "abcdefghijklmnopqrstuvwxyz0123456789..."
  }
}

provider "azurerm" {
  features {}
}
