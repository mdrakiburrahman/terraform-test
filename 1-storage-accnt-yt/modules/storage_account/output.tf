output "primary_location" {
  description = "Location where Storage Account is deployed"
  value       = azurerm_storage_account.example.primary_location
}
