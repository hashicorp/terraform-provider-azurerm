package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type RouteFilterResource struct {
}

func TestAccRouteFilter_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route_filter", "test")
	r := RouteFilterResource{}

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

func TestAccRouteFilter_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route_filter", "test")
	r := RouteFilterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_route_filter"),
		},
	})
}

func TestAccRouteFilter_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route_filter", "test")
	r := RouteFilterResource{}

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

func TestAccRouteFilter_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route_filter", "test")
	r := RouteFilterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basic,
			TestResource: r,
		}),
	})
}

func TestAccRouteFilter_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route_filter", "test")
	r := RouteFilterResource{}

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
			Config: r.withTagsUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("staging"),
			),
		},
	})
}

func TestAccRouteFilter_withRules(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route_filter", "test")
	r := RouteFilterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withRules(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("rule.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.access").HasValue("Allow"),
				check.That(data.ResourceName).Key("rule.0.rule_type").HasValue("Community"),
				check.That(data.ResourceName).Key("rule.0.communities.0").HasValue("12076:53005"),
				check.That(data.ResourceName).Key("rule.0.communities.1").HasValue("12076:53006"),
			),
		},
		{
			Config: r.withRulesUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("rule.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.access").HasValue("Allow"),
				check.That(data.ResourceName).Key("rule.0.rule_type").HasValue("Community"),
				check.That(data.ResourceName).Key("rule.0.communities.0").HasValue("12076:52005"),
				check.That(data.ResourceName).Key("rule.0.communities.1").HasValue("12076:52006"),
			),
		},
	})
}

func (t RouteFilterResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.RouteFilterID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Network.RouteFiltersClient.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		return nil, fmt.Errorf("reading Route Filter (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (RouteFilterResource) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.RouteFilterID(state.ID)
	if err != nil {
		return nil, err
	}

	future, err := client.Network.RouteFiltersClient.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("deleting Route Filter %q: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Network.RouteFiltersClient.Client); err != nil {
		return nil, fmt.Errorf("waiting for Deletion of Route Filter %q: %+v", id, err)
	}

	return utils.Bool(true), nil
}

func (RouteFilterResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_route_filter" "test" {
  name                = "acctestrf%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r RouteFilterResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_route_filter" "import" {
  name                = azurerm_route_filter.test.name
  location            = azurerm_route_filter.test.location
  resource_group_name = azurerm_route_filter.test.resource_group_name
}
`, r.basic(data))
}

func (RouteFilterResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_route_filter" "test" {
  name                = "acctestrf%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (RouteFilterResource) withTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_route_filter" "test" {
  name                = "acctestrf%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (RouteFilterResource) withTagsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_route_filter" "test" {
  name                = "acctestrf%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    environment = "staging"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (RouteFilterResource) withRules(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_route_filter" "test" {
  name                = "acctestrf%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  rule {
    name        = "acctestrule%d"
    access      = "Allow"
    rule_type   = "Community"
    communities = ["12076:53005", "12076:53006"]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (RouteFilterResource) withRulesUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_route_filter" "test" {
  name                = "acctestrf%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  rule {
    name        = "acctestrule%d"
    access      = "Allow"
    rule_type   = "Community"
    communities = ["12076:52005", "12076:52006"]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
