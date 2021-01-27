package network_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type IPGroupDataSource struct {
}

func TestAccDataSourceIPGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_ip_group", "test")
	r := IPGroupDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("cidrs.#").HasValue("0"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
	})
}

func TestAccDataSourceIpGroup_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_ip_group", "test")
	r := IPGroupDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("cidrs.#").HasValue("3"),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
			),
		},
	})
}

func (IPGroupDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_ip_group" "test" {
  name                = azurerm_ip_group.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, IPGroupResource{}.basic(data))
}

func (IPGroupDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_ip_group" "test" {
  name                = azurerm_ip_group.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, IPGroupResource{}.complete(data))
}
