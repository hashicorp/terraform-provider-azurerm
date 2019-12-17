package azurerm

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAzureRMDataSourceLoadBalancer_basic(t *testing.T) {
	dataSourceName := "data.azurerm_lb.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataSourceLoadBalancer_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "tags.%", "2"),
				),
			},
		},
	})
}

func testAccAzureRMDataSourceLoadBalancer_basic(rInt int, location string) string {
	resource := testAccAzureRMLoadBalancer_basic(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_lb" "test" {
  name                = "${azurerm_lb.test.name}"
  resource_group_name = "${azurerm_lb.test.resource_group_name}"
}
`, resource)
}
