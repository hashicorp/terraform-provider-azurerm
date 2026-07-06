# Copyright IBM Corp. 2023, 2025
# SPDX-License-Identifier: MPL-2.0

output "app_name" {
  value = azurerm_linux_function_app.example.name
}
