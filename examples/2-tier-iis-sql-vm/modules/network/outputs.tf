###################################################################################################
# Outputs
####################################################################################################

output "dc_subnet_subnet_id" {
  value = "${azurerm_subnet.dc-subnet.id}"
}

output "waf_subnet_subnet_id" {
  value = "${azurerm_subnet.waf-subnet.id}"
}

output "rp_subnet_subnet_id" {
  value = "${azurerm_subnet.rp-subnet.id}"
}

output "is_subnet_subnet_id" {
  value = "${azurerm_subnet.is-subnet.id}"
}

output "db_subnet_subnet_id" {
  value = "${azurerm_subnet.db-subnet.id}"
}

output "out_resource_group_name" {
  value = "${azurerm_resource_group.network.name}"
}
