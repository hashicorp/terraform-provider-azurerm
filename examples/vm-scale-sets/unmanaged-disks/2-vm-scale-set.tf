resource "azurerm_virtual_machine_scale_set" "example" {
  name                = "${var.prefix}-vmss"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  upgrade_policy_mode = "Manual"

  sku {
    name     = "Standard_F2"
    tier     = "Standard"
    capacity = 2
  }

  os_profile {
    computer_name_prefix = "${var.prefix}-vmss"
    admin_username       = "myadmin"
    admin_password       = "Passwword1234"
  }

  os_profile_linux_config {
    disable_password_authentication = true

    ssh_keys {
      path     = "/home/myadmin/.ssh/authorized_keys"
      key_data = "${file("~/.ssh/id_rsa.pub")}"
    }
  }

  network_profile {
    name    = "TestNetworkProfile"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      subnet_id = "${azurerm_subnet.example.id}"
    }
  }

  storage_profile_os_disk {
    name           = "osDiskProfile"
    caching        = "ReadWrite"
    create_option  = "FromImage"
    vhd_containers = ["${azurerm_storage_account.example.primary_blob_endpoint}${azurerm_storage_container.example.name}"]
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}
