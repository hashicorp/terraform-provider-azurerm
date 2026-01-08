# Copyright IBM Corp. 2014, 2025
# SPDX-License-Identifier: MPL-2.0

output "managed_instance_id" {
  value = azurerm_mssql_managed_instance.mi.id
}

output "managed_instance_name" {
  value = azurerm_mssql_managed_instance.mi.name
}

output "vnet_name" {
  value = azurerm_virtual_network.vnet.name
}