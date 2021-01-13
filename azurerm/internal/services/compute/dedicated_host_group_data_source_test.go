package compute

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type DedicatedHostGroupDataSource struct {
}

func TestAccDataSourceDedicatedHostGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_dedicated_host_group", "test")
	r := DedicatedHostGroupDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("zones.#").HasValue("1"),
				check.That(data.ResourceName).Key("zones.0").HasValue("1"),
				check.That(data.ResourceName).Key("platform_fault_domain_count").HasValue("2"),
			),
		},
	})
}

func (DedicatedHostGroupDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_dedicated_host_group" "test" {
  name                = azurerm_dedicated_host_group.test.name
  resource_group_name = azurerm_dedicated_host_group.test.resource_group_name
}
`, DedicatedHostGroupResource{}.complete(data))
}
