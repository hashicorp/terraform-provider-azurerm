# Copyright IBM Corp. 2014, 2025
# SPDX-License-Identifier: MPL-2.0

output "hostname" {
  value = var.hostname
}

output "ip_address" {
  value = azurerm_public_ip.transferpip.ip_address
}

output "fqdn" {
  value = azurerm_public_ip.transferpip.ip_address
}

output "id" {
  value = azurerm_public_ip.transferpip.id
}
