output "database_name" {
  value = "${azurerm_sql_database.db.name}"
}

output "sql_server_name" {
  value = "${azurerm_sql_server.server.name}"
}

output "sql_server_location" {
  value = "${azurerm_sql_server.server.location}"
}

output "sql_server_version" {
  value = "${azurerm_sql_server.server.version}"
}

output "sql_server_fqdn" {
  value = "${azurerm_sql_server.server.fully_qualified_domain_name}"
}

output "connection_string" {
  value = "Server=tcp:${azurerm_sql_server.server.fully_qualified_domain_name},1433;Initial Catalog=${azurerm_sql_database.db.name};Persist Security Info=False;User ID=${azurerm_sql_server.server.administrator_login};Password=${azurerm_sql_server.server.administrator_login_password};MultipleActiveResultSets=False;Encrypt=True;TrustServerCertificate=False;Connection Timeout=30;"
}
