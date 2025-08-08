# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

output "resource_group_name" {
  description = "The name of the resource group"
  value       = azurerm_resource_group.example.name
}

output "netapp_account_name" {
  description = "The name of the NetApp account"
  value       = azurerm_netapp_account.example.name
}

output "netapp_pool_name" {
  description = "The name of the NetApp capacity pool"
  value       = azurerm_netapp_pool.example.name
}

output "netapp_volume_name" {
  description = "The name of the NetApp volume"
  value       = azurerm_netapp_volume.example.name
}

output "netapp_volume_mount_ips" {
  description = "The mount IP addresses for the NetApp volume"
  value       = azurerm_netapp_volume.example.mount_ip_addresses
}

output "private_endpoint_id" {
  description = "The ID of the private endpoint to the cross-tenant key vault"
  value       = azurerm_private_endpoint.cross_tenant.id
}

output "private_endpoint_connection_status" {
  description = "The connection status of the private endpoint"
  value       = length(azurerm_private_endpoint.cross_tenant.private_service_connection) > 0 ? azurerm_private_endpoint.cross_tenant.private_service_connection[0].private_connection_resource_id : "No connection"
}

output "cross_tenant_key_vault_url" {
  description = "The constructed URL for the cross-tenant key vault key"
  value       = "https://${var.cross_tenant_key_vault_name}.vault.azure.net/keys/${var.cross_tenant_key_name}"
}

output "netapp_account_encryption_info" {
  description = "Information about the NetApp account encryption configuration"
  value = {
    encryption_key                       = azurerm_netapp_account_encryption.example.encryption_key
    user_assigned_identity_id            = azurerm_netapp_account_encryption.example.user_assigned_identity_id
    federated_client_id                  = azurerm_netapp_account_encryption.example.federated_client_id
    cross_tenant_key_vault_resource_id   = azurerm_netapp_account_encryption.example.cross_tenant_key_vault_resource_id
  }
}
