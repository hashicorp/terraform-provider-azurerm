output "azurerm_resource_group_name" {
  value = "${azurerm_resource_group.azlb.name}"
}

output "number_of_nodes" {
  value = "${var.number_of_endpoints}"
}

output "azurerm_lb_id" {
  value = "${azurerm_lb.azlb.id}"
}

output "azurerm_lb_name" {
  value = "${azurerm_lb.azlb.name}"
}

output "azurerm_lb_location" {
  value = "${azurerm_lb.azlb.location}"
}

output "azurerm_lb_frontend_ip_configuration" {
  value = "${azurerm_lb.azlb.frontend_ip_configuration}"
}

output "azurerm_lb_tags" {
  value = "${azurerm_resource_group.azlb.tags}"
}

output "azurerm_public_ip_name" {
  value = "${azurerm_public_ip.azlb.name}"
}

output "azurerm_public_ip_id" {
  value = "${azurerm_public_ip.azlb.id}"
}

output "azurerm_lb_probe_ids" {
  value = "${azurerm_lb_probe.azlb.*.id}"
}

output "azurerm_lb_nat_rule_tcp_remotevm" {
  value = "${azurerm_lb_nat_rule.azlb.*.id}"
}