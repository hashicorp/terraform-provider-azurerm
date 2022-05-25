output "app_name" {
  value = azurerm_windows_web_app.example.name
}

output "app_url" {
  value = "https://${azurerm_windows_web_app.example.default_hostname}"
}
