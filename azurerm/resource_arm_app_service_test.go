package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMAppService_basic(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMAppService_basic(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists("azurerm_app_service.test"),
				),
			},
		},
	})
}

func testCheckAzureRMAppServiceDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).appsClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_app_service" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(resourceGroup, name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("App Service still exists:\n%#v", resp)
		}
	}

	return nil
}

func testCheckAzureRMAppServiceExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		appServiceName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for App Service: %s", appServiceName)
		}

		conn := testAccProvider.Meta().(*ArmClient).appsClient

		resp, err := conn.Get(resourceGroup, appServiceName)
		if err != nil {
			return fmt.Errorf("Bad: Get on appsClient: %s", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: App Service %q (resource group: %q) does not exist", appServiceName, resourceGroup)
		}

		return nil
	}
}

func testAccAzureRMAppService_basic(rInt int) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "West Europe"
}

resource "azurerm_app_service" "test" {
    name                = "acctestAS-%d"
    location            = "West Europe"
    resource_group_name = "${azurerm_resource_group.test.name}"

		tags {
			environment = "Production"
		}
}
`, rInt, rInt)
}
