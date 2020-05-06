package tests

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceAzureRMDedicatedHostGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_dedicated_host_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDedicatedHostGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDedicatedHostGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "zones.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "zones.0", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "platform_fault_domain_count", "2"),
				),
			},
		},
	})
}

func testAccDataSourceDedicatedHostGroup_basic(data acceptance.TestData) string {
	config := testAccAzureRMDedicatedHostGroup_complete(data)
	return fmt.Sprintf(`
%s

data "azurerm_dedicated_host_group" "test" {
  name                = azurerm_dedicated_host_group.test.name
  resource_group_name = azurerm_dedicated_host_group.test.resource_group_name
}
`, config)
}
