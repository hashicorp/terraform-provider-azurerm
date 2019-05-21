package azurerm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"testing"
)

func TestAccDataSourceAzureRMRedisCache_standard(t *testing.T) {
	dataSourceName := "data.azurerm_redis_cache.test"
	ri := tf.AccRandTimeInt()

	name := fmt.Sprintf("acctestrediscache-%d", ri)
	resourceGroupName := fmt.Sprintf("acctestRG-%d", ri)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRedisCacheDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMRedisCache_standard(ri, testLocation()),
			},
			{
				Config: testAccDataSourceAzureRMRedisCache_standardWithDataSource(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", name),
					resource.TestCheckResourceAttr(dataSourceName, "resource_group_name", resourceGroupName),
					resource.TestCheckResourceAttr(dataSourceName, "ssl_port", "6380"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.environment", "test"),
				),
			},
		},
	})

}

func testAccDataSourceAzureRMRedisCache_standard(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_redis_cache" "test" {
  name                = "acctestrediscache-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  capacity            = 1
  family              = "C"
  sku_name            = "Standard"
  enable_non_ssl_port = false
  redis_configuration {}

  tags = {
    environment = "test"
  }
}
`, rInt, location, rInt)
}

func testAccDataSourceAzureRMRedisCache_standardWithDataSource(rInt int, location string) string {
	config := testAccAzureRMRedisCache_standard(rInt, testLocation())
	return fmt.Sprintf(`
%s

data "azurerm_redis_cache" "test" {
	  name                = "${azurerm_redis_cache.test.name}"
  	  resource_group_name = "${azurerm_redis_cache.test.resource_group_name}"

}
`, config)
}
