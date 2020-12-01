package trafficmanager_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type TrafficManagerProfileDataSource struct{}

func TestAccAzureRMDataSourceTrafficManagerProfile(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_traffic_manager_profile", "test")

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: TrafficManagerProfileDataSource{}.template(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("traffic_routing_method").HasValue("Performance"),
			),
		},
	})
}

func (d TrafficManagerProfileDataSource) template(data acceptance.TestData) string {
	template := TrafficManagerProfileResource{}.basic(data, "Performance")
	return fmt.Sprintf(`
%s

data "azurerm_traffic_manager_profile" "test" {
  name                = azurerm_traffic_manager_profile.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, template)
}
