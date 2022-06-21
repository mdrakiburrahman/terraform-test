output "storage_account_name" {
  value = module.storage_account.storage_account_name
}

output "blob_endpoint" {
  value = module.storage_account.primary_blob_endpoint
}

output "primary_location" {
  value = module.storage_account.primary_location
}

output "resource_group_name" {
  value = module.resource_group.resource_group_name
}
