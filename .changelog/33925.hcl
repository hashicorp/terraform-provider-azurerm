change "new-property" {
  body = "`azurerm_mssql_database` - add support for the `free_limit_enabled` and `free_limit_exhaustion_behavior` properties"
}
change "resource-fix" {
  body = "`azurerm_mssql_database` - deprecate the `Free` value of the `sku_name` property as it is no longer accepted by Azure"
}
