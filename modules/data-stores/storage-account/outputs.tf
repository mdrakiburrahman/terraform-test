output "storage_account_name" {
  description = "Primary connection endpoint for the storage account."
  value       = azurerm_storage_account.storage_accnt.name
}

output "primary_blob_endpoint" {
  description = "Primary blob connection endpoint for the storage account."
  value       = azurerm_storage_account.storage_accnt.primary_blob_endpoint
}

output "secondary_blob_endpoint" {
  description = "Secondary blob connection endpoint for the storage account."
  value       = azurerm_storage_account.storage_accnt.secondary_blob_endpoint
}

output "primary_location" {
  description = "Location where Storage Account is deployed"
  value       = azurerm_storage_account.storage_accnt.primary_location
}