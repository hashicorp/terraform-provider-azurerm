# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

output "vm_fqdn" {
  value = azurerm_public_ip.example.fqdn
}

output "ssh_command" {
  value = "ssh ${local.admin_username}@${azurerm_public_ip.example.fqdn}"
}
