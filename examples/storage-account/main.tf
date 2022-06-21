module "resource_group" {
  source              = "../../modules/misc/resource-group"
  resource_group_name = var.resource_group_name
  location            = var.location
  tags                = var.tags
}

module "storage_account" {
  depends_on               = [module.resource_group]
  source                   = "../../modules/data-stores/storage-account"
  resource_group_name      = var.resource_group_name
  location                 = var.location
  account_tier             = var.account_tier
  account_kind             = var.account_kind
  unique_postfix           = var.account_unique_postfix
  tags                     = var.tags
  account_replication_type = var.replication_type
}
