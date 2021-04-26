package network_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type RouteFilterDataSource struct {
}

func TestAccDataSourceRouteFilter_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_route_filter", "test")
	r := RouteFilterDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("route.#").HasValue("0"),
			),
		},
	})
}

func TestAccDataSourceRouteFilter_withRules(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_route_filter", "test")
	r := RouteFilterDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.withRules(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("rule.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.access").HasValue("Allow"),
				check.That(data.ResourceName).Key("rule.0.rule_type").HasValue("Community"),
				check.That(data.ResourceName).Key("rule.0.communities.0").HasValue("12076:53005"),
			),
		},
	})
}

func (RouteFilterDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_route_filter" "test" {
  name                = azurerm_route_filter.test.name
  resource_group_name = azurerm_route_filter.test.resource_group_name
}
`, RouteFilterResource{}.basic(data))
}

func (RouteFilterDataSource) withRules(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_route_filter" "test" {
  name                = azurerm_route_filter.test.name
  resource_group_name = azurerm_route_filter.test.resource_group_name
}
`, RouteFilterResource{}.withRules(data))
}
