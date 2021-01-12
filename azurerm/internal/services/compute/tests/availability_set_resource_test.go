package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type AvailabilitySetResource struct {
}

func TestAccAvailabilitySet_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_availability_set", "test")
	r := AvailabilitySetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("platform_update_domain_count").HasValue("5"),
				check.That(data.ResourceName).Key("platform_fault_domain_count").HasValue("3"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAvailabilitySet_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_availability_set", "test")
	r := AvailabilitySetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("platform_update_domain_count").HasValue("5"),
				check.That(data.ResourceName).Key("platform_fault_domain_count").HasValue("3"),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_availability_set"),
		},
	})
}

func TestAccAvailabilitySet_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_availability_set", "test")
	r := AvailabilitySetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("platform_update_domain_count").HasValue("5"),
				check.That(data.ResourceName).Key("platform_fault_domain_count").HasValue("3"),
				testCheckAvailabilitySetDisappears(data.ResourceName),
			),
			ExpectNonEmptyPlan: true,
		},
	})
}

func TestAccAvailabilitySet_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_availability_set", "test")
	r := AvailabilitySetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withTags(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("Production"),
				check.That(data.ResourceName).Key("tags.cost_center").HasValue("MSFT"),
			),
		},
		{
			Config: r.withUpdatedTags(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("staging"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAvailabilitySet_withPPG(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_availability_set", "test")
	r := AvailabilitySetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withPPG(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("proximity_placement_group_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAvailabilitySet_withDomainCounts(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_availability_set", "test")
	r := AvailabilitySetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withDomainCounts(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("platform_update_domain_count").HasValue("3"),
				check.That(data.ResourceName).Key("platform_fault_domain_count").HasValue("3"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAvailabilitySet_unmanaged(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_availability_set", "test")
	r := AvailabilitySetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.unmanaged(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("managed").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func (t AvailabilitySetResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.AvailabilitySetID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Compute.AvailabilitySetsClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Compute Availability Set %q", id.String())
	}

	return utils.Bool(resp.ID != nil), nil
}

func testCheckAvailabilitySetDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.AvailabilitySetsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		id, err := parse.AvailabilitySetID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Delete(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if !response.WasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Delete on availSetClient: %+v", err)
			}
		}

		return nil
	}
}

func (AvailabilitySetResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_availability_set" "test" {
  name                = "acctestavset-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r AvailabilitySetResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_availability_set" "import" {
  name                = azurerm_availability_set.test.name
  location            = azurerm_availability_set.test.location
  resource_group_name = azurerm_availability_set.test.resource_group_name
}
`, r.basic(data))
}

func (AvailabilitySetResource) withTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_availability_set" "test" {
  name                = "acctestavset-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (AvailabilitySetResource) withUpdatedTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_availability_set" "test" {
  name                = "acctestavset-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    environment = "staging"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (AvailabilitySetResource) withPPG(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_proximity_placement_group" "test" {
  name                = "acctestPPG-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_availability_set" "test" {
  name                = "acctestavset-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  proximity_placement_group_id = azurerm_proximity_placement_group.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (AvailabilitySetResource) withDomainCounts(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_availability_set" "test" {
  name                         = "acctestavset-%d"
  location                     = azurerm_resource_group.test.location
  resource_group_name          = azurerm_resource_group.test.name
  platform_update_domain_count = 3
  platform_fault_domain_count  = 3
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (AvailabilitySetResource) unmanaged(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_availability_set" "test" {
  name                         = "acctestavset-%d"
  location                     = azurerm_resource_group.test.location
  resource_group_name          = azurerm_resource_group.test.name
  platform_update_domain_count = 3
  platform_fault_domain_count  = 3
  managed                      = false
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
