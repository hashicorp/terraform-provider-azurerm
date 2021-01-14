provider "azurerm" {
  features {}
}

locals {
  virtual_machine_name = "${var.prefix}-vm"
}

resource "azurerm_virtual_machine" "example" {
  name                  = "${local.virtual_machine_name}"
  resource_group_name   = "${azurerm_resource_group.example.name}"
  location              = "${azurerm_resource_group.example.location}"
  network_interface_ids = ["${azurerm_network_interface.example.id}"]
  vm_size               = "Standard_F2"

  storage_image_reference {
    publisher = "radware"
    offer     = "radware-alteon-va"
    sku       = "radware-alteon-ng-va-adc"
    version   = "latest"
  }

  plan {
    name      = "radware-alteon-ng-va-adc"
    publisher = "radware"
    product   = "radware-alteon-va"
  }

  storage_os_disk {
    name              = "${local.virtual_machine_name}-osdisk"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  os_profile {
    computer_name  = "hostname"
    admin_username = "testadmin"
    admin_password = "Password1234!"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }
}
