package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccDataSourceAzureRMMapsAccount_basic(t *testing.T) {
	dataSourceName := "data.azurerm_maps_account.test"
	rInt := tf.AccRandTimeInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMMapsAccount_basic(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resource_group_name"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.environment", "testing"),
					resource.TestCheckResourceAttr(dataSourceName, "sku_name", "s0"),
					resource.TestCheckResourceAttrSet(dataSourceName, "x_ms_client_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(dataSourceName, "secondary_access_key"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMMapsAccount_basic(rInt int, location string) string {
	template := testAccAzureRMMapsAccount_tags(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_maps_account" "test" {
  name                = azurerm_maps_account.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, template)
}
