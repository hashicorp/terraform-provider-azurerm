output "resource_group_name" {
  value = data.azurerm_resource_group.rg.name
}

output "resource_group_location" {
  value = data.azurerm_resource_group.rg.location
}

output "function_app_name" {
  value = azurerm_linux_function_app.example.name
}

output "frontend_url" {
  value = "${azurerm_api_management.example.gateway_url}/${azurerm_api_management_api.example.path}${azurerm_api_management_api_operation.example.url_template}"
}
