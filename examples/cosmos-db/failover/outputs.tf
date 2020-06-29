output "cosmos-db-id" {
  value = "${azurerm_cosmosdb_account.example.id}"
}

output "cosmos-db-endpoint" {
  value = "${azurerm_cosmosdb_account.example.endpoint}"
}

output "cosmos-db-endpoints_read" {
  value = "${azurerm_cosmosdb_account.example.read_endpoints}"
}

output "cosmos-db-endpoints_write" {
  value = "${azurerm_cosmosdb_account.example.write_endpoints}"
}

output "cosmos-db-primary_master_key" {
  value = "${azurerm_cosmosdb_account.example.primary_master_key}"
}

output "cosmos-db-secondary_master_key" {
  value = "${azurerm_cosmosdb_account.example.secondary_master_key}"
}
