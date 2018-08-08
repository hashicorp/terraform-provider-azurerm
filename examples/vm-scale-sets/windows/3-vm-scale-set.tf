resource "azurerm_virtual_machine_scale_set" "main" {
  name                = "${local.virtual_machine_name}"
  location            = "${azurerm_resource_group.main.location}"
  resource_group_name = "${azurerm_resource_group.main.name}"
  upgrade_policy_mode = "Manual"
  overprovision       = true

  sku {
    name     = "Standard_F2"
    tier     = "Standard"
    capacity = "${var.instance_count}"
  }

  os_profile {
    computer_name_prefix = "${local.virtual_machine_name}"
    admin_username       = "${var.admin_username}"
    admin_password       = "${var.admin_password}"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  network_profile {
    name    = "${local.virtual_machine_name}-nic"
    primary = true

    ip_configuration {
      name                                   = "${local.virtual_machine_name}-ipconfig"
      subnet_id                              = "${azurerm_subnet.internal.id}"
      load_balancer_backend_address_pool_ids = ["${azurerm_lb_backend_address_pool.main.id}"]
      load_balancer_inbound_nat_rules_ids    = ["${element(azurerm_lb_nat_pool.main.*.id, count.index)}"]
    }
  }

  storage_profile_os_disk {
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  storage_profile_data_disk {
    lun           = 0
    caching       = "ReadWrite"
    create_option = "Empty"
    disk_size_gb  = 10
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04.0-LTS"
    version   = "latest"
  }
}
