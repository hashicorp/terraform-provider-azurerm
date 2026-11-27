// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_mover_job_definition_test

import (
	"fmt"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

type StorageMoverAgentResource struct{}

func (r StorageMoverAgentResource) template(data acceptance.TestData) string {
	clientData := data.Client()
	randomUUID, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestnw-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_network_interface" "test" {
  name                = "acctestnic-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.test.id
  }
}

resource "azurerm_network_security_group" "my_terraform_nsg" {
  name                = "myNetworkSecurityGroup"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

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
  network_interface_id      = azurerm_network_interface.test.id
  network_security_group_id = azurerm_network_security_group.my_terraform_nsg.id
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
}

resource "azurerm_linux_virtual_machine" "test" {
  name                            = "acctestVM-%[1]d"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  size                            = "Standard_F2"
  admin_username                  = "adminuser"
  admin_password                  = "TerraformTest01!"
  provision_vm_agent              = false
  allow_extension_operations      = false
  disable_password_authentication = false
  network_interface_ids = [
    azurerm_network_interface.test.id,
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

  connection {
    type     = "ssh"
    host     = azurerm_public_ip.test.ip_address
    user     = "adminuser"
    password = "TerraformTest01!"
  }

  provisioner "file" {
    content = templatefile("scripts/install_arc.sh.tftpl", {
      resource_group_name = azurerm_resource_group.test.name
      uuid                = "%[3]s"
      location            = azurerm_resource_group.test.location
      tenant_id           = data.azurerm_client_config.current.tenant_id
      client_id           = "%[4]s"
      client_secret       = "%[5]s"
      subscription_id     = data.azurerm_client_config.current.subscription_id
    })
    destination = "/home/adminuser/install_arc_agent.sh"
  }

  provisioner "remote-exec" {
    inline = [
      "sudo apt-get install -y python-ctypes",
      "sudo sed -i 's/\r$//' /home/adminuser/install_arc_agent.sh",
      "sudo chmod +x /home/adminuser/install_arc_agent.sh",
      "bash /home/adminuser/install_arc_agent.sh",
    ]
  }
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest-uai-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_resource_group.test.id
  role_definition_name = "Contributor"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_storage_mover" "test" {
  name                = "acctest-ssm-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  depends_on = [
    azurerm_linux_virtual_machine.test
  ]
}

data "azurerm_arc_machine" "test" {
  name                = azurerm_linux_virtual_machine.test.name
  resource_group_name = azurerm_resource_group.test.name
  depends_on = [
    azurerm_linux_virtual_machine.test
  ]
}


`, data.RandomInteger, data.Locations.Primary, randomUUID, clientData.Default.ClientID, clientData.Default.ClientSecret)
}
