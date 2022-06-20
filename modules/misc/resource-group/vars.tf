variable "resource_group_name" {
  type = string
}

variable "location" {
  type    = string
  default = "eastus"
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
