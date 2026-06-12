# Copyright IBM Corp. 2014, 2025
# SPDX-License-Identifier: MPL-2.0

output "sql_server_fqdn" {
  value = azurerm_mssql_server.example.fully_qualified_domain_name
}

output "database_name" {
  value = azurerm_mssql_database.example.name
}
