resource "azurerm_virtual_machine_scale_set" "example" {
  name                = "${var.prefix}-vmss"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  upgrade_policy_mode = "Manual"

  sku {
    name     = "Standard_D1_v2"
    tier     = "Standard"
    capacity = 2
  }

  os_profile {
    computer_name_prefix = "${var.prefix}-vm"
    admin_username       = "myadmin"
    admin_password       = "Passwword1234"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  network_profile {
    name    = "web_ss_net_profile"
    primary = true

    ip_configuration {
      name                                   = "internal"
      subnet_id                              = "${azurerm_subnet.example.id}"
      load_balancer_backend_address_pool_ids = ["${azurerm_lb_backend_address_pool.example.id}"]
      primary                                = true
    }
  }

  storage_profile_os_disk {
    name              = ""
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  extension {
    name                 = "CustomScriptForLinux"
    publisher            = "Microsoft.OSTCExtensions"
    type                 = "CustomScriptForLinux"
    type_handler_version = "1.4"

    settings = <<SETTINGS
    {
      "commandToExecute" : "sudo apt-get -y install apache2"
    }
SETTINGS
  }
}
