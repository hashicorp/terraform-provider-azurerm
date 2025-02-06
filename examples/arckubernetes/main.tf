# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

data "azurerm_subscription" "current" {}
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-rg"
  location = var.location
}

resource "azurerm_virtual_network" "example" {
  name                = "${var.prefix}-vn"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name                 = "${var.prefix}-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_public_ip" "example" {
  name                = "${var.prefix}-pi"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  allocation_method   = "Static"
}

resource "azurerm_network_interface" "example" {
  name                = "${var.prefix}-ni"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.example.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.example.id
  }
}

resource "azurerm_network_security_group" "example" {
  name                = "${var.prefix}NetworkSecurityGroup"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  security_rule {
    name                       = "SSH"
    priority                   = 1001
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "22"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
  }
}

resource "azurerm_network_interface_security_group_association" "example" {
  network_interface_id      = azurerm_network_interface.example.id
  network_security_group_id = azurerm_network_security_group.example.id
}

resource "azurerm_linux_virtual_machine" "example" {
  name                            = "${var.prefix}-lvm"
  resource_group_name             = azurerm_resource_group.example.name
  location                        = azurerm_resource_group.example.location
  size                            = "Standard_F2"
  admin_username                  = var.user_name
  admin_password                  = var.password
  provision_vm_agent              = false
  allow_extension_operations      = false
  disable_password_authentication = false
  network_interface_ids = [
    azurerm_network_interface.example.id,
  ]
  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }
  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }
}


resource "azurerm_arc_kubernetes_cluster" "example" {
  name                         = "${var.prefix}-akc"
  resource_group_name          = azurerm_resource_group.example.name
  location                     = azurerm_resource_group.example.location
  agent_public_key_certificate = var.public_key
  identity {
    type = "SystemAssigned"
  }


  connection {
    type     = "ssh"
    host     = azurerm_public_ip.example.ip_address
    user     = var.user_name
    password = var.password
  }

  provisioner "file" {
    content = templatefile("testdata/install_agent.sh.tftpl", {
      subscription_id     = data.azurerm_subscription.current.subscription_id
      resource_group_name = azurerm_resource_group.example.name
      cluster_name        = azurerm_arc_kubernetes_cluster.example.name
      location            = azurerm_resource_group.example.location
      tenant_id           = data.azurerm_client_config.current.tenant_id
      working_dir         = "/home/${var.user_name}"
    })
    destination = "/home/${var.user_name}/install_agent.sh"
  }

  provisioner "file" {
    source      = "testdata/install_agent.py"
    destination = "/home/${var.user_name}/install_agent.py"
  }

  provisioner "file" {
    source      = "testdata/kind.yaml"
    destination = "/home/${var.user_name}/kind.yaml"
  }

  provisioner "file" {
    content     = var.private_pem
    destination = "/home/${var.user_name}/private.pem"
  }

  provisioner "remote-exec" {
    inline = [
      "sudo sed -i 's/\r$//' /home/${var.user_name}/install_agent.sh",
      "sudo chmod +x /home/${var.user_name}/install_agent.sh",
      "bash /home/${var.user_name}/install_agent.sh > /home/${var.user_name}/agent_log",
    ]
  }


  depends_on = [
    azurerm_linux_virtual_machine.example
  ]
}
