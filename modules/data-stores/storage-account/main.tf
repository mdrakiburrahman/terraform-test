resource "azurerm_storage_account" "storage_accnt" {
  name                     = lower("storageacct${var.unique_postfix}")
  resource_group_name      = var.resource_group_name
  location                 = var.location
  account_tier             = var.account_tier
  account_kind             = var.account_kind
  account_replication_type = var.account_replication_type

  tags = var.tags
}
