package compute_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type DedicatedHostDataSource struct {
}

func TestAccDataSourceAzureRMDedicatedHost_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_dedicated_host", "test")
	r := DedicatedHostDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("tags.%").Exists(),
			),
		},
	})
}

func (DedicatedHostDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_dedicated_host" "test" {
  name                      = azurerm_dedicated_host.test.name
  dedicated_host_group_name = azurerm_dedicated_host_group.test.name
  resource_group_name       = azurerm_dedicated_host_group.test.resource_group_name
}
`, DedicatedHostResource{}.basic(data))
}
