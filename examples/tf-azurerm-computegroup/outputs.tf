output "vmss_name" { 
    value = "${azurerm_virtual_machine_scale_set.vmss.name}" 
}

output "vnet_id" { 
    value = "${azurerm_virtual_network.vnet.id}"
}

output "subnet_id" { 
    value = "${azurerm_subnet.subnet1.id}"
}

