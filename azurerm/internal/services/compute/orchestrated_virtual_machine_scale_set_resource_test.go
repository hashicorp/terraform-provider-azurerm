package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMOrchestratedVirtualMachineScaleSet_basicZonal(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orchestrated_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMOrchestratedVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMOrchestratedVirtualMachineScaleSet_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMOrchestratedVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMOrchestratedVirtualMachineScaleSet_updateZonal(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orchestrated_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMOrchestratedVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMOrchestratedVirtualMachineScaleSet_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMOrchestratedVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMOrchestratedVirtualMachineScaleSet_update(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMOrchestratedVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMOrchestratedVirtualMachineScaleSet_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMOrchestratedVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMOrchestratedVirtualMachineScaleSet_basicNonZonal(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orchestrated_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMOrchestratedVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMOrchestratedVirtualMachineScaleSet_basicNonZonal(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMOrchestratedVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMOrchestratedVirtualMachineScaleSet_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orchestrated_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMOrchestratedVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMOrchestratedVirtualMachineScaleSet_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMOrchestratedVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMOrchestratedVirtualMachineScaleSet_requiresImport),
		},
	})
}

func testCheckAzureRMOrchestratedVirtualMachineScaleSetDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.VMScaleSetClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_orchestrated_virtual_machine_scale_set" {
			continue
		}

		id, err := parse.VirtualMachineScaleSetID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Get on Compute.VMScaleSetClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testCheckAzureRMOrchestratedVirtualMachineScaleSetExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.VMScaleSetClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		id, err := parse.VirtualMachineScaleSetID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Virtual Machine Scale Set VM Mode %q (Resource Group: %q) does not exist", id.Name, id.ResourceGroup)
			}

			return fmt.Errorf("bad: Get on Compute.VMScaleSetClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMOrchestratedVirtualMachineScaleSet_basic(data acceptance.TestData) string {
	template := testAccAzureRMOrchestratedVirtualMachineScaleSet_template(data)
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
`, template, data.RandomInteger)
}

func testAccAzureRMOrchestratedVirtualMachineScaleSet_update(data acceptance.TestData) string {
	template := testAccAzureRMOrchestratedVirtualMachineScaleSet_template(data)
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
`, template, data.RandomInteger)
}

func testAccAzureRMOrchestratedVirtualMachineScaleSet_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMOrchestratedVirtualMachineScaleSet_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_orchestrated_virtual_machine_scale_set" "import" {
  name                = azurerm_orchestrated_virtual_machine_scale_set.test.name
  location            = azurerm_orchestrated_virtual_machine_scale_set.test.location
  resource_group_name = azurerm_orchestrated_virtual_machine_scale_set.test.resource_group_name

  platform_fault_domain_count = azurerm_orchestrated_virtual_machine_scale_set.test.platform_fault_domain_count
  single_placement_group      = azurerm_orchestrated_virtual_machine_scale_set.test.single_placement_group
}
`, template)
}

func testAccAzureRMOrchestratedVirtualMachineScaleSet_basicNonZonal(data acceptance.TestData) string {
	template := testAccAzureRMOrchestratedVirtualMachineScaleSet_template(data)
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
`, template, data.RandomInteger)
}

func testAccAzureRMOrchestratedVirtualMachineScaleSet_template(data acceptance.TestData) string {
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
