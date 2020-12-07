package netapp_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func testAccDataSourceAzureRMNetAppAccount_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_netapp_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNetAppAccount_basicConfig(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
				),
			},
		},
	})
}

func testAccDataSourceNetAppAccount_basicConfig(data acceptance.TestData) string {
	config := testAccAzureRMNetAppAccount_basicConfig(data)
	return fmt.Sprintf(`
%s

data "azurerm_netapp_account" "test" {
  resource_group_name = azurerm_netapp_account.test.resource_group_name
  name                = azurerm_netapp_account.test.name
}
`, config)
}
