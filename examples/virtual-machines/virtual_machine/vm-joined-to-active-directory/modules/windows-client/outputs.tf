# Copyright IBM Corp. 2014, 2025
# SPDX-License-Identifier: MPL-2.0

output "public_ip_address" {
  value = azurerm_public_ip.static.ip_address
}
