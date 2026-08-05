# Copyright IBM Corp. 2014, 2025
# SPDX-License-Identifier: MPL-2.0

output "app_name" {
  value = azurerm_linux_web_app.example.name
}

output "app_url" {
  value = "https://${azurerm_linux_web_app.example.default_hostname}"
}
