output "azurerm_vnet_id" {
  value = "${azurerm_virtual_network.vnet.id}"
}

output "azurerm_vnet_name" {
  value = "${azurerm_virtual_network.vnet.name}"
}

output "azurerm_vnet_location" {
  value = "${azurerm_virtual_network.vnet.location}"
}

output "azurerm_vnet_address_space" {
  value = "${azurerm_virtual_network.vnet.location}"
}

output "azurerm_vnet_dns_servers" {
  value = "${azurerm_virtual_network.vnet.dns_servers}"
}

output "azurerm_vnet_subnet" {
  value = "${azurerm_virtual_network.vnet.dns_servers}"
}
