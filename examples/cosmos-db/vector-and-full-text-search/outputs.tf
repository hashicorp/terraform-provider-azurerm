# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

output "cosmosdb-account-id" {
  value = azurerm_cosmosdb_account.example.id
}

output "cosmosdb-account-endpoint" {
  value = azurerm_cosmosdb_account.example.endpoint
}

output "cosmosdb-account-primary_key" {
  value = azurerm_cosmosdb_account.example.primary_key
}

output "cosmosdb-database-id" {
  value = azurerm_cosmosdb_sql_database.example.id
}

output "cosmosdb-container-id" {
  value = azurerm_cosmosdb_sql_container.example.id
}
