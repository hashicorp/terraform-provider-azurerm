# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

output "id" {
  value = azurerm_virtual_machine.example.id
}

output "network_interface_id" {
  value = azurerm_network_interface.example.id
}
