package maintenance_test

import (
	`context`
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/maintenance/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type MaintenanceAssignmentDedicatedHostResource struct {
}

func TestAccMaintenanceAssignmentDedicatedHost_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_maintenance_assignment_dedicated_host", "test")
	r := MaintenanceAssignmentDedicatedHostResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("location"),
	})
}

func TestAccMaintenanceAssignmentDedicatedHost_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_maintenance_assignment_dedicated_host", "test")
	r := MaintenanceAssignmentDedicatedHostResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (MaintenanceAssignmentDedicatedHostResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.MaintenanceAssignmentDedicatedHostID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Maintenance.ConfigurationAssignmentsClient.ListParent(ctx, id.DedicatedHostId.ResourceGroup, "Microsoft.Compute", "hostGroups", id.DedicatedHostId.HostGroupName, "hosts", id.DedicatedHostId.HostName)
	if err != nil {
		return nil, fmt.Errorf("retrieving Maintenance Assignment Dedicated Host (target resource id: %q): %v", id.DedicatedHostIdRaw, err)
	}

	return utils.Bool(resp.Value == nil || len(*resp.Value) == 0), nil
}


func (r MaintenanceAssignmentDedicatedHostResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_maintenance_assignment_dedicated_host" "test" {
  location                     = azurerm_resource_group.test.location
  maintenance_configuration_id = azurerm_maintenance_configuration.test.id
  dedicated_host_id            = azurerm_dedicated_host.test.id
}
`, r.template(data))
}

func (r MaintenanceAssignmentDedicatedHostResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_maintenance_assignment_dedicated_host" "import" {
  location                     = azurerm_maintenance_assignment_dedicated_host.test.location
  maintenance_configuration_id = azurerm_maintenance_assignment_dedicated_host.test.maintenance_configuration_id
  dedicated_host_id            = azurerm_maintenance_assignment_dedicated_host.test.dedicated_host_id
}
`, r.basic(data))
}

func (MaintenanceAssignmentDedicatedHostResource) template(data acceptance.TestData) string {
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
