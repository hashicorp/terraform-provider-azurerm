resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = "${var.location}"
}

resource "azurerm_storage_account" "example" {
  name                     = "${var.prefix}sa"
  resource_group_name      = "${azurerm_resource_group.example.name}"
  location                 = "${azurerm_resource_group.example.location}"
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_redis_cache" "example" {
  # NOTE: the Name used needs to be globally unique
  name                = "${var.prefix}-redis"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  capacity            = 3
  family              = "P"
  sku_name            = "Premium"
  enable_non_ssl_port = false

  redis_configuration {
    rdb_backup_enabled            = true
    rdb_backup_frequency          = 60
    rdb_backup_max_snapshot_count = 1
    rdb_storage_connection_string = "DefaultEndpointsProtocol=https;BlobEndpoint=${azurerm_storage_account.example.primary_blob_endpoint};AccountName=${azurerm_storage_account.example.name};AccountKey=${azurerm_storage_account.example.primary_access_key}"
  }
}
