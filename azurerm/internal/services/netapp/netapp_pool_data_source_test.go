package netapp_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMNetAppPool_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_netapp_pool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNetAppPool_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "account_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "service_level"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "size_in_tb"),
				),
			},
		},
	})
}

func testAccDataSourceNetAppPool_basic(data acceptance.TestData) string {
	config := testAccAzureRMNetAppPool_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_netapp_pool" "test" {
  resource_group_name = azurerm_netapp_pool.test.resource_group_name
  account_name        = azurerm_netapp_pool.test.account_name
  name                = azurerm_netapp_pool.test.name
}
`, config)
}
