package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMRouteFilter_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_route_filter", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRouteFilterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMRouteFilter_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteFilterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "route.#", "0"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMRouteFilter_withRules(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_route_filter", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRouteFilterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMRouteFilter_withRules(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteFilterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.access", "Allow"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.rule_type", "Community"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.communities.0", "12076:53005"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMRouteFilter_basic(data acceptance.TestData) string {
	r := testAccAzureRMRouteFilter_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_route_filter" "test" {
  name                = azurerm_route_filter.test.name
  resource_group_name = azurerm_route_filter.test.resource_group_name
}
`, r)
}

func testAccDataSourceAzureRMRouteFilter_withRules(data acceptance.TestData) string {
	r := testAccAzureRMRouteFilter_withRules(data)
	return fmt.Sprintf(`
%s

data "azurerm_route_filter" "test" {
  name                = azurerm_route_filter.test.name
  resource_group_name = azurerm_route_filter.test.resource_group_name
}
`, r)
}
