package automanage_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automanage/2022-05-04/configurationprofilehcrpassignments"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ArcMachineConfigurationAssignmentResource struct{}

func TestAccArcMachineConfigurationAssignment_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_arc_machine_automanage_configuration_assignment", "test")
	r := ArcMachineConfigurationAssignmentResource{}
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

func TestAccArcMachineConfigurationAssignment_requireImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_arc_machine_automanage_configuration_assignment", "test")
	r := ArcMachineConfigurationAssignmentResource{}
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

func (r ArcMachineConfigurationAssignmentResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	client := clients.Automanage.ConfigurationProfileArcMachineAssignmentsClient

	id, err := configurationprofilehcrpassignments.ParseProviders2ConfigurationProfileAssignmentID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func (r ArcMachineConfigurationAssignmentResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_arc_machine_automanage_configuration_assignment" "import" {
  arc_machine_id   = azurerm_arc_machine_automanage_configuration_assignment.test.arc_machine_id
  configuration_id = azurerm_arc_machine_automanage_configuration_assignment.test.configuration_id
}
`, config)
}

func (r ArcMachineConfigurationAssignmentResource) basic(data acceptance.TestData) string {
	secret := os.Getenv("ARM_CLIENT_SECRET")
	randomUUID, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

data "azurerm_client_config" "current" {}

# note: real-life usage prefer random_uuid resource in registry.terraform.io/hashicorp/random
locals {
  random_uuid = "%s"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestnw-%d"
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
  name                = "acctestnic-%d"
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
  name                = "acctestpip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
}

resource "azurerm_linux_virtual_machine" "test" {
  name                            = "acctestVM-%d"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  size                            = "Standard_F2"
  admin_username                  = "adminuser"
  admin_password                  = "AdminPassword0123!"
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
    password = "AdminPassword0123!"
  }

  provisioner "file" {
    content = templatefile("scripts/install_arc.sh.tftpl", {
      resource_group_name = azurerm_resource_group.test.name
      uuid                = local.random_uuid
      location            = azurerm_resource_group.test.location
      tenant_id           = data.azurerm_client_config.current.tenant_id
      client_id           = data.azurerm_client_config.current.client_id
      client_secret       = "%s"
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

data "azurerm_arc_machine" "test" {
  name                = azurerm_linux_virtual_machine.test.name
  resource_group_name = azurerm_resource_group.test.name
  depends_on = [
    azurerm_linux_virtual_machine.test
  ]
}

resource "azurerm_automanage_configuration" "test" {
  name                = "acctest-amcp-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_arc_machine_automanage_configuration_assignment" "test" {
  arc_machine_id   = data.azurerm_arc_machine.test.id
  configuration_id = azurerm_automanage_configuration.test.id

  depends_on = [
    azurerm_automanage_configuration.test,
    azurerm_linux_virtual_machine.test
  ]
}
`, randomUUID, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, secret, data.RandomInteger)
}
