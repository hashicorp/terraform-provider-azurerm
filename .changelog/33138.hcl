change "resource-enhancement" {
  body = "`azurerm_mysql_flexible_database` - suppress persistent diffs on `charset` and `collation` when using the `utf8` aliases as the API now returns `utf8mb3` prefixed values"
}
