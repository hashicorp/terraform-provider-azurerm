package compute_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type OrchestratedVirtualMachineScaleSetResource struct {
}

func TestAccOrchestratedVirtualMachineScaleSet_basicZonal(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orchestrated_virtual_machine_scale_set", "test")
	r := OrchestratedVirtualMachineScaleSetResource{}

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

func TestAccOrchestratedVirtualMachineScaleSet_updateZonal(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orchestrated_virtual_machine_scale_set", "test")
	r := OrchestratedVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
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
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccOrchestratedVirtualMachineScaleSet_basicNonZonal(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orchestrated_virtual_machine_scale_set", "test")
	r := OrchestratedVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicNonZonal(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccOrchestratedVirtualMachineScaleSet_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orchestrated_virtual_machine_scale_set", "test")
	r := OrchestratedVirtualMachineScaleSetResource{}

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

func (t OrchestratedVirtualMachineScaleSetResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.VirtualMachineScaleSetID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Compute.VMScaleSetClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Compute Marketplace Agreement %q", id.String())
	}

	return utils.Bool(resp.ID != nil), nil
}

func (r OrchestratedVirtualMachineScaleSetResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_orchestrated_virtual_machine_scale_set" "test" {
  name                = "acctestVMO-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  platform_fault_domain_count = 1

  zones = ["1"]

  tags = {
    ENV = "Test"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r OrchestratedVirtualMachineScaleSetResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_orchestrated_virtual_machine_scale_set" "test" {
  name                = "acctestVMO-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  platform_fault_domain_count = 1

  zones = ["1"]

  tags = {
    ENV = "Test",
    FOO = "Bar"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r OrchestratedVirtualMachineScaleSetResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_orchestrated_virtual_machine_scale_set" "import" {
  name                = azurerm_orchestrated_virtual_machine_scale_set.test.name
  location            = azurerm_orchestrated_virtual_machine_scale_set.test.location
  resource_group_name = azurerm_orchestrated_virtual_machine_scale_set.test.resource_group_name

  platform_fault_domain_count = azurerm_orchestrated_virtual_machine_scale_set.test.platform_fault_domain_count
  single_placement_group      = azurerm_orchestrated_virtual_machine_scale_set.test.single_placement_group
}
`, r.basic(data))
}

func (r OrchestratedVirtualMachineScaleSetResource) basicNonZonal(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_orchestrated_virtual_machine_scale_set" "test" {
  name                = "acctestVMO-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  platform_fault_domain_count = 2

  tags = {
    ENV = "Test"
  }
}
`, r.template(data), data.RandomInteger)
}

func (OrchestratedVirtualMachineScaleSetResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-VMSS-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}
