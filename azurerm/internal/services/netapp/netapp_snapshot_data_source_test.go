package netapp_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMNetAppSnapshot_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_netapp_snapshot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNetAppSnapshot_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
				),
			},
		},
	})
}

func testAccDataSourceNetAppSnapshot_basic(data acceptance.TestData) string {
	config := testAccAzureRMNetAppSnapshot_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_netapp_snapshot" "test" {
  resource_group_name = azurerm_netapp_snapshot.test.resource_group_name
  account_name        = azurerm_netapp_snapshot.test.account_name
  pool_name           = azurerm_netapp_snapshot.test.pool_name
  volume_name         = azurerm_netapp_snapshot.test.volume_name
  name                = azurerm_netapp_snapshot.test.name
}
`, config)
}
