terraform {
  required_version = "~> 1.0"
  required_providers {
    azurerm = "~> 3.9.0"
  }
}

provider "azurerm" {
  features {}
}
