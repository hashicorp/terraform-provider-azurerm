###################################################################################################
# Outputs
####################################################################################################

output "out_resource_group_name" {
  value = "${azurerm_network_interface.primary.resource_group_name}"
}

output "out_dc_location" {
  value = "${azurerm_virtual_machine_extension.create-active-directory-forest.location}"
}

output "out_dcavailabilityset" {
  value = "${azurerm_availability_set.dcavailabilityset.id}"
}
