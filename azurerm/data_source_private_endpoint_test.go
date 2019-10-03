package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccDataSourceAzureRMPrivateEndpoint_basic(t *testing.T) {
	dataSourceName := "data.azurerm_private_endpoint.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePrivateEndpoint_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "subnet_id"),
					resource.TestCheckResourceAttr(dataSourceName, "private_link_service_connections.0.name", fmt.Sprintf("acctestconnection-%d", ri)),
					resource.TestCheckResourceAttrSet(dataSourceName, "private_link_service_connections.0.private_link_service_id"),
					resource.TestCheckResourceAttr(dataSourceName, "private_link_service_connections.0.request_message", "Please approve my request"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMPrivateEndpoint_complete(t *testing.T) {
	dataSourceName := "data.azurerm_private_endpoint.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePrivateEndpoint_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "subnet_id"),
					resource.TestCheckResourceAttr(dataSourceName, "private_link_service_connections.0.name", fmt.Sprintf("acctestconnection-%d", ri)),
					resource.TestCheckResourceAttrSet(dataSourceName, "private_link_service_connections.0.private_link_service_id"),
					resource.TestCheckResourceAttr(dataSourceName, "private_link_service_connections.0.group_ids.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "private_link_service_connections.0.request_message", "plz approve my request"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.env", "test"),
				),
			},
		},
	})
}

func testAccDataSourcePrivateEndpoint_basic(rInt int, location string) string {
	config := testAccAzureRMPrivateEndpoint_basic(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_private_endpoint" "test" {
  resource_group_name = "${azurerm_private_endpoint.test.resource_group_name}"
  name                = "${azurerm_private_endpoint.test.name}"
}
`, config)
}

func testAccDataSourcePrivateEndpoint_complete(rInt int, location string) string {
	config := testAccAzureRMPrivateEndpoint_complete(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_private_endpoint" "test" {
  resource_group_name = "${azurerm_private_endpoint.test.resource_group_name}"
  name                = "${azurerm_private_endpoint.test.name}"
}
`, config)
}
