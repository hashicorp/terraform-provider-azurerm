package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccDataSourceAzureRMSpringCloud_complete(t *testing.T) {
	dataSourceName := "data.azurerm_spring_cloud.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSpringCloud_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "service_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "version"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.env", "test"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.version", "1"),
				),
			},
		},
	})
}

func testAccDataSourceSpringCloud_complete(rInt int, location string) string {
	config := testAccAzureRMSpringCloud_complete(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_spring_cloud" "test" {
  name                = azurerm_spring_cloud.test.name
  resource_group      = azurerm_spring_cloud.test.resource_group
}
`, config)
}
