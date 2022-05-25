output "app_name" {
  value = azurerm_windows_web_app.example.name
}

output "app_service_default_hostname" {
  value = "https://${azurerm_windows_web_app.example.default_hostname}"
}
