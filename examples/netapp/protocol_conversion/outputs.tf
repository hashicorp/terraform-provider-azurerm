output "volume_id" {
  description = "The ID of the NetApp volume"
  value       = azurerm_netapp_volume.example.id
}

output "volume_mount_ip_addresses" {
  description = "The mount IP addresses for the NetApp volume"
  value       = azurerm_netapp_volume.example.mount_ip_addresses
}

output "current_protocol_type" {
  description = "The current protocol type of the volume"
  value       = var.protocol_type
}

output "resource_group_name" {
  description = "The name of the resource group"
  value       = azurerm_resource_group.example.name
}