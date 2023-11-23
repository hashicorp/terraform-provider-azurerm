package systemcentervirtualmachinemanager_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/vmmservers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SystemCenterVirtualMachineManagerServerResource struct{}

func TestAccSystemCenterVirtualMachineManagerServer_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_system_center_virtual_machine_manager_server", "test")
	r := SystemCenterVirtualMachineManagerServerResource{}

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

func TestAccSystemCenterVirtualMachineManagerServer_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_system_center_virtual_machine_manager_server", "test")
	r := SystemCenterVirtualMachineManagerServerResource{}

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

func TestAccSystemCenterVirtualMachineManagerServer_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_system_center_virtual_machine_manager_server", "test")
	r := SystemCenterVirtualMachineManagerServerResource{}

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

func TestAccSystemCenterVirtualMachineManagerServer_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_system_center_virtual_machine_manager_server", "test")
	r := SystemCenterVirtualMachineManagerServerResource{}

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

func (r SystemCenterVirtualMachineManagerServerResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := vmmservers.ParseVMmServerID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.SystemCenterVirtualMachineManager.VMmServers.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r SystemCenterVirtualMachineManagerServerResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_system_center_virtual_machine_manager_server" "test" {
  name                = "acctest-scvmmms-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = azurerm_custom_location.test.id
  fqdn                = "testdomain.com"
}
`, r.template(data), data.RandomInteger)
}

func (r SystemCenterVirtualMachineManagerServerResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_system_center_virtual_machine_manager_server" "import" {
  name                = azurerm_system_center_virtual_machine_manager_server.test.name
  resource_group_name = azurerm_system_center_virtual_machine_manager_server.test.resource_group_name
  location            = azurerm_system_center_virtual_machine_manager_server.test.location
  custom_location_id  = azurerm_system_center_virtual_machine_manager_server.test.custom_location_id
  fqdn                = azurerm_system_center_virtual_machine_manager_server.test.fqdn
}
`, r.basic(data))
}

func (r SystemCenterVirtualMachineManagerServerResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_system_center_virtual_machine_manager_server" "test" {
  name                = "acctest-scvmmms-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = azurerm_custom_location.test.id
  fqdn                = "testdomain.com"
  port                = 10000

  credential {
    username = "adminTerraform"
    password = "QAZwsx123"
  }

  tags = {
    Env = "Test"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r SystemCenterVirtualMachineManagerServerResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_system_center_virtual_machine_manager_server" "test" {
  name                = "acctest-scvmmms-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = azurerm_custom_location.test.id
  fqdn                = "testdomain2.com"
  port                = 10001

  credential {
    username = "adminTerraform2"
    password = "QAZwsx124"
  }

  tags = {
    Env = "Test2"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r SystemCenterVirtualMachineManagerServerResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-scvmmms-%d"
  location = "%s"
}

resource "azurerm_custom_location" "test" {
  name = "acctest-cl-%d"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
