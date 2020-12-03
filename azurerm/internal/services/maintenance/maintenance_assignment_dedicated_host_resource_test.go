package maintenance_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/maintenance/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccMaintenanceAssignmentDedicatedHost_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_maintenance_assignment_dedicated_host", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckMaintenanceAssignmentDedicatedHostDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMaintenanceAssignmentDedicatedHost_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckMaintenanceAssignmentDedicatedHostExists(data.ResourceName),
				),
			},
			data.ImportStep("location"),
		},
	})
}

func TestAccMaintenanceAssignmentDedicatedHost_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_maintenance_assignment_dedicated_host", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckMaintenanceAssignmentDedicatedHostDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMaintenanceAssignmentDedicatedHost_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckMaintenanceAssignmentDedicatedHostExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccMaintenanceAssignmentDedicatedHost_requiresImport),
		},
	})
}

func testCheckMaintenanceAssignmentDedicatedHostDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).Maintenance.ConfigurationAssignmentsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_maintenance_assignment_dedicated_host" {
			continue
		}

		id, err := parse.MaintenanceAssignmentDedicatedHostID(rs.Primary.ID)
		if err != nil {
			return err
		}

		listResp, err := conn.ListParent(ctx, id.DedicatedHostId.ResourceGroup, "Microsoft.Compute", "hostGroups", id.DedicatedHostId.HostGroupName, "hosts", id.DedicatedHostId.HostName)
		if err != nil {
			if !utils.ResponseWasNotFound(listResp.Response) {
				return err
			}
			return nil
		}
		if listResp.Value != nil && len(*listResp.Value) > 0 {
			return fmt.Errorf("maintenance assignment (Dedicated Host ID: %q) still exists", id.DedicatedHostIdRaw)
		}

		return nil
	}

	return nil
}

func testCheckMaintenanceAssignmentDedicatedHostExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).Maintenance.ConfigurationAssignmentsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		id, err := parse.MaintenanceAssignmentDedicatedHostID(rs.Primary.ID)
		if err != nil {
			return err
		}

		listResp, err := conn.ListParent(ctx, id.DedicatedHostId.ResourceGroup, "Microsoft.Compute", "hostGroups", id.DedicatedHostId.HostGroupName, "hosts", id.DedicatedHostId.HostName)
		if err != nil {
			return fmt.Errorf("bad: list on ConfigurationAssignmentsClient: %+v", err)
		}
		if listResp.Value == nil || len(*listResp.Value) == 0 {
			return fmt.Errorf("could not find Maintenance Assignment (target resource id: %q)", id.DedicatedHostIdRaw)
		}

		return nil
	}
}

func testAccMaintenanceAssignmentDedicatedHost_basic(data acceptance.TestData) string {
	template := testAccMaintenanceAssignmentDedicatedHost_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_maintenance_assignment_dedicated_host" "test" {
  location                     = azurerm_resource_group.test.location
  maintenance_configuration_id = azurerm_maintenance_configuration.test.id
  dedicated_host_id            = azurerm_dedicated_host.test.id
}
`, template)
}

func testAccMaintenanceAssignmentDedicatedHost_requiresImport(data acceptance.TestData) string {
	template := testAccMaintenanceAssignmentDedicatedHost_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_maintenance_assignment_dedicated_host" "import" {
  location                     = azurerm_maintenance_assignment_dedicated_host.test.location
  maintenance_configuration_id = azurerm_maintenance_assignment_dedicated_host.test.maintenance_configuration_id
  dedicated_host_id            = azurerm_maintenance_assignment_dedicated_host.test.dedicated_host_id
}
`, template)
}

func testAccMaintenanceAssignmentDedicatedHost_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-maint-%[1]d"
  location = "%[2]s"
}

resource "azurerm_maintenance_configuration" "test" {
  name                = "acctest-MC%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  scope               = "All"
}

resource "azurerm_dedicated_host_group" "test" {
  name                        = "acctest-DHG-%[1]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  platform_fault_domain_count = 2
}

resource "azurerm_dedicated_host" "test" {
  name                    = "acctest-DH-%[1]d"
  location                = azurerm_resource_group.test.location
  dedicated_host_group_id = azurerm_dedicated_host_group.test.id
  sku_name                = "DSv3-Type1"
  platform_fault_domain   = 1
}
`, data.RandomInteger, data.Locations.Primary)
}
