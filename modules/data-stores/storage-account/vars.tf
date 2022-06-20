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

variable "account_kind" {
  type    = string
  default = "StorageV2"
}

variable "account_replication_type" {
  type    = string
  default = "LRS"
}

variable "unique_postfix" {
  type = string
}

variable "tags" {
  type        = map(string)
  description = "A map of the tags to use on the resources that are deployed with this module."

  default = {
    Source  = "terraform"
    Owner   = "Raki Rahman"
    Project = "Terraform Hackathon"
  }
}
