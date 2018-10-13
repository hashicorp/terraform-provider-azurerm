package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceAzureRMApiManagementApi_basic(t *testing.T) {
	dataSourceName := "data.azurerm_api_management_api.test"
	rInt := acctest.RandInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceApiManagementApi_basic(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "revision", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "is_current", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "path", "api1"),
				),
			},
		},
	})
}

func testAccDataSourceApiManagementApi_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "amtestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name            = "acctestAM-%d"
  publisher_name  = "pub1"
  publisher_email = "pub1@email.com"

  sku {
    name     = "Developer"
    capacity = 1
  }

  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_api_management_api" "test" {
  name            = "acctestAMA-%d"
	service_name    = "${azurerm_api_management.test.name}"
	path 						= "api1"

	import {
    content_value = "https://api.swaggerhub.com/apis/sparebanken-vest/tf-simple/1.0.2"
    content_format = "swagger-link-json"
  }

  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

data "azurerm_api_management_api" "test" {
  name                = "${azurerm_api_management_api.test.name}"
	service_name    		= "${azurerm_api_management.test.name}"
  resource_group_name = "${azurerm_api_management_api.test.resource_group_name}"
}
`, rInt, location, rInt, rInt)
}
