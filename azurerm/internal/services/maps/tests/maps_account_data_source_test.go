package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMMapsAccount_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_maps_account", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMMapsAccount_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "resource_group_name"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "testing"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku_name", "S0"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "x_ms_client_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_access_key"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMMapsAccount_basic(data acceptance.TestData) string {
	template := testAccAzureRMMapsAccount_tags(data)
	return fmt.Sprintf(`
%s

data "azurerm_maps_account" "test" {
  name                = azurerm_maps_account.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, template)
}
