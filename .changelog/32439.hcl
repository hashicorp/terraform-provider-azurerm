change "resource-enhancement" {
  body = "`azurerm_cosmosdb_account` - Add restriction to prevent `azurerm_cosmosdb_account` creation/update with `analytical_storage_enabled` = true for `kind` = \"GlobalDocumentDB\" as this is no longer supported"
}