output "account_id" {
  value = azurerm_function_app.main.identity.0.principal_id
}