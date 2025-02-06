// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storagemover_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2023-03-01/agents"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type StorageMoverAgentTestResource struct{}

func TestAccStorageMoverAgent_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_mover_agent", "test")
	r := StorageMoverAgentTestResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageMoverAgent_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_mover_agent", "test")
	r := StorageMoverAgentTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccStorageMoverAgent_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_mover_agent", "test")
	r := StorageMoverAgentTestResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageMoverAgent_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_mover_agent", "test")
	r := StorageMoverAgentTestResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r StorageMoverAgentTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := agents.ParseAgentID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.StorageMover.AgentsClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r StorageMoverAgentTestResource) template(data acceptance.TestData) string {
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
      client_id           = data.azurerm_client_config.current.client_id
      client_secret       = "%[4]s"
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


`, data.RandomInteger, data.Locations.Primary, randomUUID, os.Getenv("ARM_CLIENT_SECRET"))
}

func (r StorageMoverAgentTestResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`


provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

%s

resource "azurerm_storage_mover_agent" "test" {
  name                     = "acctest-sa-%d"
  storage_mover_id         = azurerm_storage_mover.test.id
  arc_virtual_machine_id   = data.azurerm_arc_machine.test.id
  arc_virtual_machine_uuid = data.azurerm_arc_machine.test.vm_uuid
}
`, template, data.RandomInteger)
}

func (r StorageMoverAgentTestResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_mover_agent" "import" {
  name                     = azurerm_storage_mover_agent.test.name
  storage_mover_id         = azurerm_storage_mover_agent.test.storage_mover_id
  arc_virtual_machine_id   = azurerm_storage_mover_agent.test.arc_virtual_machine_id
  arc_virtual_machine_uuid = azurerm_storage_mover_agent.test.arc_virtual_machine_uuid
}


`, config)
}

func (r StorageMoverAgentTestResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`


provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

%s

resource "azurerm_storage_mover_agent" "test" {
  name                     = "acctest-sa-%d"
  storage_mover_id         = azurerm_storage_mover.test.id
  arc_virtual_machine_id   = data.azurerm_arc_machine.test.id
  arc_virtual_machine_uuid = data.azurerm_arc_machine.test.vm_uuid
  description              = "Example Agent Description"
}
`, template, data.RandomInteger)
}

func (r StorageMoverAgentTestResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`


provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

%s

resource "azurerm_storage_mover_agent" "test" {
  name                     = "acctest-sa-%d"
  storage_mover_id         = azurerm_storage_mover.test.id
  arc_virtual_machine_id   = data.azurerm_arc_machine.test.id
  arc_virtual_machine_uuid = data.azurerm_arc_machine.test.vm_uuid
  description              = "Update Example Agent Description"

}
`, template, data.RandomInteger)
}
