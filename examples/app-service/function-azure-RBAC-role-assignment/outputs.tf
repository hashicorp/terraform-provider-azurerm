# Copyright IBM Corp. 2014, 2025
# SPDX-License-Identifier: MPL-2.0

output "account_id" {
  value = azurerm_linux_function_app.example.identity.0.principal_id
}