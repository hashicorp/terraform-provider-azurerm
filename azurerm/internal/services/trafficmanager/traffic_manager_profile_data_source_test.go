package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccAzureRMDataSourceTrafficManagerProfile(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_traffic_manager_profile", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataSourceTrafficManagerProfile_template(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "traffic_routing_method", "Performance"),
				),
			},
		},
	})
}

func testAccAzureRMDataSourceTrafficManagerProfile_template(data acceptance.TestData) string {
	template := testAccAzureRMTrafficManagerProfile_basic(data, "Performance")
	return fmt.Sprintf(`
%s

data "azurerm_traffic_manager_profile" "test" {
  name                = azurerm_traffic_manager_profile.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, template)
}
