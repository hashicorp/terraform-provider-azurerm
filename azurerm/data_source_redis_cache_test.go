package azurerm

import (
	"fmt"

	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccDataSourceAzureRMRedisCache_standard(t *testing.T) {
	dataSourceName := "data.azurerm_redis_cache.test"
	ri := tf.AccRandTimeInt()

	name := fmt.Sprintf("acctestRedis-%d", ri)
	resourceGroupName := fmt.Sprintf("acctestRG-%d", ri)
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRedisCacheDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMRedisCache_standardWithDataSource(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", name),
					resource.TestCheckResourceAttr(dataSourceName, "resource_group_name", resourceGroupName),
					resource.TestCheckResourceAttr(dataSourceName, "ssl_port", "6380"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.environment", "production"),
				),
			},
		},
	})

}

func testAccDataSourceAzureRMRedisCache_standardWithDataSource(rInt int, location string) string {
	config := testAccAzureRMRedisCache_standard(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_redis_cache" "test" {
  name                = "${azurerm_redis_cache.test.name}"
  resource_group_name = "${azurerm_redis_cache.test.resource_group_name}"
}
`, config)
}
