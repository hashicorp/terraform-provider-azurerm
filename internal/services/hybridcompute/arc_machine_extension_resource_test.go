// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package hybridcompute_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2022-11-10/machineextensions"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ArcMachineExtensionResource struct {
}

func TestAccArcMachineExtension_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_arc_machine_extension", "test")
	r := ArcMachineExtensionResource{}
	template := r.template(data)
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, template),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("publisher").HasValue("Microsoft.Azure.Monitor"),
				check.That(data.ResourceName).Key("type").HasValue("AzureMonitorLinuxAgent"),
			),
		},
		data.ImportStep("protected_settings"),
	})
}

func TestAccArcMachineExtension_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_arc_machine_extension", "test")
	r := ArcMachineExtensionResource{}
	template := r.template(data)
	basicConfig := r.basic(data, template)
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: basicConfig,
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(basicConfig),
			ExpectError: acceptance.RequiresImportError("azurerm_arc_machine_extension"),
		},
	})
}

func TestAccArcMachineExtension_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_arc_machine_extension", "test")
	r := ArcMachineExtensionResource{}
	template := r.template(data)
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, template),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("publisher").HasValue("Microsoft.Azure.Extensions"),
				check.That(data.ResourceName).Key("type").HasValue("CustomScript"),
				check.That(data.ResourceName).Key("type_handler_version").MatchesRegex(regexp.MustCompile("^2[.]1.*$")),
				check.That(data.ResourceName).Key("automatic_upgrade_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("settings").HasValue(`{"timestamp":123456789}`),
			),
		},
		data.ImportStep("protected_settings"),
	})
}

func TestAccArcMachineExtension_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_arc_machine_extension", "test")
	r := ArcMachineExtensionResource{}
	template := r.template(data)
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, template),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("protected_settings"),
		{
			Config: r.update(data, template),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("protected_settings"),
	})
}

func (r ArcMachineExtensionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := machineextensions.ParseExtensionID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.HybridCompute.MachineExtensionsClient
	resp, err := client.Get(ctx, *id)
	exists := false
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return &exists, nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func (r ArcMachineExtensionResource) basic(data acceptance.TestData, template string) string {
	return fmt.Sprintf(`
				%s

resource "azurerm_arc_machine_extension" "test" {
  name           = "acctest-hcme-%d"
  arc_machine_id = data.azurerm_arc_machine.test.id
  publisher      = "Microsoft.Azure.Monitor"
  type           = "AzureMonitorLinuxAgent"
  location       = "%s"
  lifecycle {
    ignore_changes = [type_handler_version]
  }
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r ArcMachineExtensionResource) requiresImport(basicConfig string) string {
	return fmt.Sprintf(`
			%s

resource "azurerm_arc_machine_extension" "import" {
  name           = azurerm_arc_machine_extension.test.name
  arc_machine_id = azurerm_arc_machine_extension.test.arc_machine_id
  publisher      = azurerm_arc_machine_extension.test.publisher
  type           = azurerm_arc_machine_extension.test.type
  location       = azurerm_arc_machine_extension.test.location
  lifecycle {
    ignore_changes = [type_handler_version]
  }
}
`, basicConfig)
}

func (r ArcMachineExtensionResource) complete(data acceptance.TestData, template string) string {
	return fmt.Sprintf(`
			%s

resource "azurerm_arc_machine_extension" "test" {
  name                      = "acctest-hcme-%d"
  arc_machine_id            = data.azurerm_arc_machine.test.id
  location                  = "%s"
  automatic_upgrade_enabled = false
  publisher                 = "Microsoft.Azure.Extensions"
  settings                  = jsonencode({ "timestamp" : 123456789 })
  protected_settings        = jsonencode({ "commandToExecute" : "echo 'Hello World!'" })
  type                      = "CustomScript"
  type_handler_version      = "2.1"

  tags = {
    Environment = "Production"
  }
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r ArcMachineExtensionResource) update(data acceptance.TestData, template string) string {
	return fmt.Sprintf(`
			%s

resource "azurerm_arc_machine_extension" "test" {
  name                      = "acctest-hcme-%d"
  arc_machine_id            = data.azurerm_arc_machine.test.id
  location                  = "%s"
  automatic_upgrade_enabled = true
  publisher                 = "Microsoft.Azure.Monitor"
  type                      = "AzureMonitorLinuxAgent"
  type_handler_version      = "1.24"
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r ArcMachineExtensionResource) template(data acceptance.TestData) string {
	secret := os.Getenv("ARM_CLIENT_SECRET")
	randomUUID, _ := uuid.GenerateUUID()
	password := generateRandomPassword(10)
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
  admin_password                  = "%s"
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
    offer     = "UbuntuServer"
    sku       = "18.04-LTS"
    version   = "latest"
  }

  connection {
    type     = "ssh"
    host     = azurerm_public_ip.test.ip_address
    user     = "adminuser"
    password = "%s"
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
`, randomUUID, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, password, password, secret)
}
