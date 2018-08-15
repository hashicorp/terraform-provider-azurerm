output "cosmos-db-id" {
  value = "${azurerm_cosmosdb_account.db.id}"
}

output "cosmos-db-endpoint" {
  value = "${azurerm_cosmosdb_account.db.endpoint}"
}

output "cosmos-db-endpoints_read" {
  value = "${azurerm_cosmosdb_account.db.read_endpoints}"
}

output "cosmos-db-endpoints_write" {
  value = "${azurerm_cosmosdb_account.db.write_endpoints}"
}

output "cosmos-db-primary_master_key" {
  value = "${azurerm_cosmosdb_account.db.primary_master_key}"
}

output "cosmos-db-secondary_master_key" {
  value = "${azurerm_cosmosdb_account.db.secondary_master_key}"
}
