# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

output "app_name" {
  value = azurerm_linux_function_app.example.name
}

output "function_url" {
  value = "${azurerm_function_app_function.example.invocation_url}?"
}
