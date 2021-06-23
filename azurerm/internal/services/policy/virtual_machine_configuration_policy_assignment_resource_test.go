package policy_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/policy/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type VirtualMachineConfigurationPolicyAssignmentResource struct{}

func TestAccVirtualMachineConfigurationPolicyAssignment_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_configuration_policy_assignment", "test")
	r := VirtualMachineConfigurationPolicyAssignmentResource{}

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

func TestAccVirtualMachineConfigurationPolicyAssignment_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_configuration_policy_assignment", "test")
	r := VirtualMachineConfigurationPolicyAssignmentResource{}

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

func (r VirtualMachineConfigurationPolicyAssignmentResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.VirtualMachineConfigurationPolicyAssignmentID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Policy.GuestConfigurationAssignmentsClient.Get(ctx, id.ResourceGroup, id.GuestConfigurationAssignmentName, id.VirtualMachineName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Properties != nil), nil
}

func (r VirtualMachineConfigurationPolicyAssignmentResource) templateBase(data acceptance.TestData) string {
	return fmt.Sprintf(`
locals {
  vm_name = "acctestvm%s"
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
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
  name                = "acctestnic-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}
`, data.RandomString, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r VirtualMachineConfigurationPolicyAssignmentResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }
}
`, r.templateBase(data))
}

func (r VirtualMachineConfigurationPolicyAssignmentResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_configuration_policy_assignment" "test" {
  name               = "acctest-gca-%d"
  location           = azurerm_windows_virtual_machine.test.location
  virtual_machine_id = azurerm_windows_virtual_machine.test.id
  configuration {
    name    = "WhitelistedApplication"
    version = "1.*"

    parameter {
      name  = "[InstalledApplication]bwhitelistedapp;Name"
      value = "NotePad,sql"
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r VirtualMachineConfigurationPolicyAssignmentResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_configuration_policy_assignment" "import" {
  name               = azurerm_virtual_machine_configuration_policy_assignment.test.name
  location           = azurerm_virtual_machine_configuration_policy_assignment.test.location
  virtual_machine_id = azurerm_virtual_machine_configuration_policy_assignment.test.virtual_machine_id

  configuration {
    name    = "WhitelistedApplication"
    version = "1.*"

    parameter {
      name  = "[InstalledApplication]bwhitelistedapp;Name"
      value = "NotePad,sql"
    }
  }
}
`, r.basic(data))
}
