# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

output "openshift_console_url" {
  value = "https://${azurerm_public_ip.openshift_primary_pip.fqdn}:8443/console"
}

output "openshift_primary_ssh" {
  value = "ssh ${var.admin_username}@${azurerm_public_ip.openshift_primary_pip.fqdn} -p 2200"
}

output "openshift_infra_load_balancer_fqdn" {
  value = azurerm_public_ip.infra_lb_pip.fqdn
}

output "node_os_storage_account_name" {
  value = azurerm_storage_account.nodeos_storage_account.name
}

output "node_data_storage_account_name" {
  value = azurerm_storage_account.nodedata_storage_account.name
}

output "infra_storage_account_name" {
  value = azurerm_storage_account.infra_storage_account.name
}
