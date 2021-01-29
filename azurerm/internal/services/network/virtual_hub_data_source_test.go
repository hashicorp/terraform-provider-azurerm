package network_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type VirtualHubDataSource struct {
}

func TestAccDataSourceAzureRMVirtualHub_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_virtual_hub", "test")
	r := VirtualHubDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("address_prefix").Exists(),
				check.That(data.ResourceName).Key("virtual_wan_id").Exists(),
			),
		},
	})
}

func (VirtualHubDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_virtual_hub" "test" {
  name                = azurerm_virtual_hub.test.name
  resource_group_name = azurerm_virtual_hub.test.resource_group_name
}
`, VirtualHubResource{}.basic(data))
}
