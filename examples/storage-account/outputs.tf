output "blob_endpoint" {
  value = module.storage_account.primary_blob_endpoint
}

output "resource_group_name" {
  value = module.resource_group.resource_group_name
}
