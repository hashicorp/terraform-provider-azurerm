output "azurerm_public_ip_name" {
  value = "${azurerm_public_ip.mypublicIP.name}"
}

output "azurerm_lb_id" {
  value = "${azurerm_lb.mylb.id}"
}

output "azurerm_lb_private_ip_address" {
  value = "${azurerm_lb.mylb.private_ip_address}"
}