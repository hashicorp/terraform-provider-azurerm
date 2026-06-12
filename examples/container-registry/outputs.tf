# Copyright IBM Corp. 2014, 2025
# SPDX-License-Identifier: MPL-2.0

output "login_server" {
  value = azurerm_container_registry.example.login_server
}
