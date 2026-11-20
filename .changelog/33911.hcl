change "resource-fix" {
  body = "`azurerm_mssql_database` - prevent sending `max_size_gb` when updating a database with `create_mode` set to `Secondary` or `OnlineSecondary`"
}
