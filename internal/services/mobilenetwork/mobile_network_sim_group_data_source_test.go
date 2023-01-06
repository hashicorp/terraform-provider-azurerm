package mobilenetwork_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type MobileNetworkSimGroupDataSource struct{}

func TestAccMobileNetworkSimGroupDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_sim_group", "test")
	d := MobileNetworkSimGroupDataSource{}
	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key(`location`).HasValue(data.Locations.Primary),
				check.That(data.ResourceName).Key(`mobile_network_id`).Exists(),
				check.That(data.ResourceName).Key(`encryption_key_url`).Exists(),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func (r MobileNetworkSimGroupDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%s

data "azurerm_mobile_network_sim_group" "test" {
  name                = azurerm_mobile_network_sim_group.test.name
  resource_group_name = azurerm_mobile_network_sim_group.test.resource_group_name
}
`, MobileNetworkSimGroupResource{}.complete(data))
}
