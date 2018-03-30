
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

output "cosmos-db-geo_location-0-id" {
  value = "${azurerm_cosmosdb_account.db.geo_location}"
}

/*output "cosmos-db-geo_location-1-endpoint" {
  value = "${azurerm_cosmosdb_account.db.geo_location.0.document_endpoint}"
}*/

output "cosmos-db-primary_master_key" {
  value = "${azurerm_cosmosdb_account.db.primary_master_key}"
}

output "cosmos-db-secondary_master_key" {
  value = "${azurerm_cosmosdb_account.db.secondary_master_key}"
}