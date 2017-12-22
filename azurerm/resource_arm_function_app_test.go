package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMFunctionApp_basic(t *testing.T) {
	resourceName := "azurerm_function_app.test"
	ri := acctest.RandInt()
	config := testAccAzureRMFunctionApp_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(resourceName),
				),
			},
		},
	})
}

func testCheckAzureRMFunctionAppDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).appServicesClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_function_app" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(resourceGroup, name)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}
			return err
		}

		return fmt.Errorf("Function App still exists:\n%#v", resp)
	}

	return nil
}

func testCheckAzureRMFunctionAppExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		functionAppName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Function App: %s", functionAppName)
		}

		client := testAccProvider.Meta().(*ArmClient).appServicesClient

		resp, err := client.Get(resourceGroup, functionAppName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Function App %q (resource group: %q) does not exist", functionAppName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on appServicesClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMFunctionApp_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
	name     = "acctestRG-%d"
	location = "%s"
}

resource "azurerm_storage_account" "test" {
	name                     = "acctestsa%d"
	resource_group_name      = "${azurerm_resource_group.test.name}"
	location                 = "${azurerm_resource_group.test.location}"
	account_tier             = "Standard"
	account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
	name                = "acctestASP-%d"
	location            = "${azurerm_resource_group.test.location}"
	resource_group_name = "${azurerm_resource_group.test.name}"

	sku {
		tier = "Standard"
		size = "S1"
	}
}

resource "azurerm_function_app" "test" {
	name                      = "acctestFA-%d"
	location                  = "${azurerm_resource_group.test.location}"
	resource_group_name       = "${azurerm_resource_group.test.name}"
	app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
	storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"
}
`, rInt, location, rInt%1000000000000, rInt, rInt)
}
