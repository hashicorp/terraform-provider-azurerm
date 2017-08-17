output "azurerm_lb_id" {
  value = "${azurerm_lb.mylb.id}"
}

output "azurerm_lb_name" {
  value = "${azurerm_lb.mylb.name}"
}

output "azurerm_lb_location" {
  value = "${azurerm_lb.mylb.location}"
}

output "azurerm_lb_frontend_ip_configuration" {
  value = "${azurerm_lb.mylb.frontend_ip_configuration}"
}

output "azurerm_lb_tags" {
  value = "${azurerm_lb.mylb.tags}"
}

output "azurerm_public_ip_name" {
  value = "${azurerm_public_ip.mypublicIP.name}"
}

output "vnet_id" {
  value = "${azurerm_virtual_network.vnet.id}"
}

output "vnet_name" {
  value = "${azurerm_virtual_network.vnet.name}"
}

output "vnet_location" {
  value = "${azurerm_virtual_network.vnet.location}"
}

output "vnet_address_space" {
  value = "${azurerm_virtual_network.vnet.address_space}"
}

output "vnet_dns_servers" {
  value = "${azurerm_virtual_network.vnet.dns_servers}"
}

output "vnet_subnet" {
  value = "${azurerm_virtual_network.vnet.subnet}"
}

output "azurerm_network_interface_ids" {
  value = "${azurerm_network_interface.nic.*.id}"
}

output "azurerm_lb_probe_ids" {
  value = "${azurerm_lb_probe.lb_probe.*.id}"
}

output "azurerm_lb_nat_rule_tcp-remotevm" {
  value = "${azurerm_lb_nat_rule.tcp-remotevm.*.*}"
}