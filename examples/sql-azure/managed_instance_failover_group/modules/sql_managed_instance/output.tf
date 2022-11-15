output "managed_instance_id" {
  value = azurerm_sql_managed_instance.mi.id
}

output "managed_instance_name" {
  value = azurerm_sql_managed_instance.mi.name
}

output "vnet_name" {
  value = azurerm_virtual_network.vnet.name
}