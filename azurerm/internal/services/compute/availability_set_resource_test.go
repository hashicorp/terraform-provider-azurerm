package compute_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type AvailabilitySetResource struct {
}

func TestAccAvailabilitySet_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_availability_set", "test")
	r := AvailabilitySetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basic,
			TestResource: r,
		}),
	})
}

func TestAccAvailabilitySet_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_availability_set", "test")
	r := AvailabilitySetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withTags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("Production"),
				check.That(data.ResourceName).Key("tags.cost_center").HasValue("MSFT"),
			),
		},
		{
			Config: r.withUpdatedTags(data),
			Check: acceptance.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withPPG(data),
			Check: acceptance.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withDomainCounts(data),
			Check: acceptance.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.unmanaged(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("managed").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func (AvailabilitySetResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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

func (AvailabilitySetResource) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.AvailabilitySetID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Compute.AvailabilitySetsClient.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if !response.WasNotFound(resp.Response) {
			return nil, fmt.Errorf("deleting on availSetClient: %+v", err)
		}
	}

	return utils.Bool(true), nil
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
