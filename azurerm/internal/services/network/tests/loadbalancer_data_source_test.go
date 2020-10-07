package tests

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAzureRMDataSourceLoadBalancer_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_lb", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataSourceLoadBalancer_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "2"),
				),
			},
		},
	})
}

func testAccAzureRMDataSourceLoadBalancer_basic(data acceptance.TestData) string {
	resource := testAccAzureRMLoadBalancer_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_lb" "test" {
  name                = azurerm_lb.test.name
  resource_group_name = azurerm_lb.test.resource_group_name
}
`, resource)
}
