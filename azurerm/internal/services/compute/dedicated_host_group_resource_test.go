package compute

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type DedicatedHostGroupResource struct {
}

func TestAccDedicatedHostGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dedicated_host_group", "test")
	r := DedicatedHostGroupResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDedicatedHostGroup_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dedicated_host_group", "test")
	r := DedicatedHostGroupResource{}

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

func TestAccDedicatedHostGroup_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dedicated_host_group", "test")
	r := DedicatedHostGroupResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("zones.#").HasValue("1"),
				check.That(data.ResourceName).Key("zones.0").HasValue("1"),
				check.That(data.ResourceName).Key("platform_fault_domain_count").HasValue("2"),
				check.That(data.ResourceName).Key("tags.ENV").HasValue("prod"),
			),
		},
		data.ImportStep(),
	})
}

func (t DedicatedHostGroupResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	resourceGroupName := id.ResourceGroup
	name := id.Path["hostGroups"]

	resp, err := clients.Compute.DedicatedHostGroupsClient.Get(ctx, resourceGroupName, name, "")
	if err != nil {
		return nil, fmt.Errorf("retrieving Compute Dedicated Host Group %q", id)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (DedicatedHostGroupResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-compute-%d"
  location = "%s"
}

resource "azurerm_dedicated_host_group" "test" {
  name                        = "acctestDHG-compute-%d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  platform_fault_domain_count = 2
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r DedicatedHostGroupResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_dedicated_host_group" "import" {
  resource_group_name         = azurerm_dedicated_host_group.test.resource_group_name
  name                        = azurerm_dedicated_host_group.test.name
  location                    = azurerm_dedicated_host_group.test.location
  platform_fault_domain_count = 2
}
`, r.basic(data))
}

func (DedicatedHostGroupResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-compute-%d"
  location = "%s"
}

resource "azurerm_dedicated_host_group" "test" {
  name                        = "acctestDHG-compute-%d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  platform_fault_domain_count = 2
  zones                       = ["1"]
  tags = {
    ENV = "prod"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
