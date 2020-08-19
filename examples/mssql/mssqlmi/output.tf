output "Managed_Instance" {
  value = azurerm_mssql_managed_instance.example
}

output "Managed_Instance_Source" {
  value = data.azurerm_mssql_managed_instance.example
}

