# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

output "login_server" {
  value = azurerm_container_registry.example.login_server
}
