variable "resource_group_name" {
  type = string
}

variable "location" {
  type    = string
  default = "eastus"
}

variable "account_tier" {
  type = string
}

variable "unique_postfix" {
  type = string
}

variable "account_kind" {
  type    = string
  default = "StorageV2"
}
