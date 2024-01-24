# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

output "app_name" {
  value = azurerm_linux_function_app.example.name
}

output "app_url" {
  value = "https://${azurerm_linux_function_app.example.default_hostname}"
}

output "scm_type" {
  value = azurerm_app_service_source_control.example.scm_type
}
