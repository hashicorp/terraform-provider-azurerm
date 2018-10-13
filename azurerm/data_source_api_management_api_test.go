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
					resource.TestCheckResourceAttr(dataSourceName, "api_id", "api1"),
					resource.TestCheckResourceAttr(dataSourceName, "import.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "import.0.content_format", "swagger-json"),
					resource.TestCheckResourceAttr(dataSourceName, "import.0.wsdl_selector.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "is_current", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "is_online", "false"),
					resource.TestCheckResourceAttr(dataSourceName, "name", "api1"),
					resource.TestCheckResourceAttr(dataSourceName, "path", "api1"),
					resource.TestCheckResourceAttr(dataSourceName, "protocols.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "protocols.0", "https"),
					resource.TestCheckResourceAttr(dataSourceName, "service_name", fmt.Sprintf("acctestAM-%d", rInt)),
					resource.TestCheckResourceAttr(dataSourceName, "soap_pass_through", "false"),
					resource.TestCheckResourceAttr(dataSourceName, "subscription_key_parameter_names.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "subscription_key_parameter_names.0.header", "Ocp-Apim-Subscription-Key"),
					resource.TestCheckResourceAttr(dataSourceName, "subscription_key_parameter_names.0.query", "subscription-key"),
					resource.TestCheckResourceAttr(dataSourceName, "version", ""),
					resource.TestCheckResourceAttr(dataSourceName, "version_set_id", ""),
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
    content_value  = "https://api.swaggerhub.com/apis/sparebanken-vest/tf-simple/1.0.2"
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
