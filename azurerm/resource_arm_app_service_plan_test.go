package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMAppServicePlan_basic(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMAppServicePlan_basic(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServicePlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServicePlanExists("azurerm_app_service_plan.test"),
				),
			},
		},
	})
}

func testCheckAzureRMAppServicePlanDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).appServicePlansClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_app_service_plan" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(resourceGroup, name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("App Service Plan still exists:\n%#v", resp)
		}
	}

	return nil
}

func testCheckAzureRMAppServicePlanExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		appServicePlanName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for App Service Plan: %s", appServicePlanName)
		}

		conn := testAccProvider.Meta().(*ArmClient).appServicePlansClient

		resp, err := conn.Get(resourceGroup, appServicePlanName)
		if err != nil {
			return fmt.Errorf("Bad: Get on appServicePlansClient: %s", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: App Service Plan %q (resource group: %q) does not exist", appServicePlanName, resourceGroup)
		}

		return nil
	}
}

func testAccAzureRMAppServicePlan_basic(rInt int) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "West Europe"
}
resource "azurerm_app_service_plan" "test" {
    name                = "acctestAppServicePlan-%d"
    location            = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    sku {
			tier = "Standard"
			size = "S0"
		}
}
`, rInt, rInt)
}
