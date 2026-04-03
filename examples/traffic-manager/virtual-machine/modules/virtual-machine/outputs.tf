# Copyright IBM Corp. 2014, 2025
# SPDX-License-Identifier: MPL-2.0

output "id" {
  value = azurerm_virtual_machine.example.id
}

output "network_interface_id" {
  value = azurerm_network_interface.example.id
}
