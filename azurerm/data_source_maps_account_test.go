package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccDataSourceAzureRMMapsAccount(t *testing.T) {
	dataSourceName := "data.azurerm_maps_account.test"
	rInt := tf.AccRandTimeInt()
	location := testLocation()
	key := "environment"
	value := "testing"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMMapsAccount(rInt, location, key, value),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resource_group_name"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.environment", value),
					resource.TestCheckResourceAttr(dataSourceName, "sku", "s0"),
					resource.TestCheckResourceAttrSet(dataSourceName, "x_ms_client_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(dataSourceName, "secondary_access_key"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMMapsAccount(rInt int, location string, key string, value string) string {
	template := testAccAzureRMMapsAccount_tags(rInt, location, key, value)
	return fmt.Sprintf(`
%s

data "azurerm_maps_account" "test" {
    name = azurerm_maps_account.test.name
    resource_group_name = azurerm_resource_group.test.name
}
`, template)
}
