change "resource-fix" {
  body = "`azurerm_mssql_server` - allow `administrator_login_password` and `administrator_login_password_wo` to be updated in the same apply that sets `azuread_administrator.0.azuread_authentication_only` to `false`"
}
