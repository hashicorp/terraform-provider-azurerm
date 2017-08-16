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