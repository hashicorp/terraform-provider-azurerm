provider "azurerm" {
  features {}
}

locals {
  virtual_machine_name = "${var.prefix}-vm"
  admin_username       = "testadmin"
}

resource "azurerm_network_interface" "example" {
  name                      = "${azurerm_resource_group.example.name}-nic"
  location                  = "${azurerm_resource_group.example.location}"
  resource_group_name       = "${azurerm_resource_group.example.name}"
  network_security_group_id = "${azurerm_network_security_group.bastion.id}"

  ip_configuration {
    name                          = "internal"
    subnet_id                     = "${azurerm_subnet.bastion.id}"
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = "${azurerm_public_ip.example.id}"
  }
}

resource "azurerm_public_ip" "example" {
  name                = "${var.prefix}-bastionpip"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  allocation_method   = "Dynamic"
}

resource "azurerm_virtual_machine" "example" {
  name                  = "${local.virtual_machine_name}"
  location              = "${azurerm_resource_group.example.location}"
  resource_group_name   = "${azurerm_resource_group.example.name}"
  network_interface_ids = ["${azurerm_network_interface.example.id}"]
  vm_size               = "Standard_F2"

  storage_image_reference {
    publisher = "MicrosoftOSTC"
    offer     = "FreeBSD"
    sku       = "11.1"
    version   = "latest"
  }

  storage_os_disk {
    name              = "${local.virtual_machine_name}-osdisk"
    managed_disk_type = "Standard_LRS"
    caching           = "ReadWrite"
    create_option     = "FromImage"
  }

  os_profile {
    computer_name  = "${local.virtual_machine_name}"
    admin_username = "${local.admin_username}"
  }

  os_profile_linux_config {
    disable_password_authentication = true

    ssh_keys {
      path     = "/home/${local.admin_username}/.ssh/authorized_keys"
      key_data = "${local.public_ssh_key}"
    }
  }
}
