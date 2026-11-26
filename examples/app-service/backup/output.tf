# Copyright IBM Corp. 2014, 2026
# SPDX-License-Identifier: MPL-2.0

output "app_name" {
  value = azurerm_windows_web_app.example.default_hostname
}

output "app_url" {
  value = "https://${azurerm_windows_web_app.example.default_hostname}"
}
