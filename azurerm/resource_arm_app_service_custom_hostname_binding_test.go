package azurerm

import (
	"fmt"
	"testing"

	"os"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func testAccAzureRMAppServiceCustomHostnameBinding_basic(t *testing.T) {
	appServiceEnvVariable := "ARM_TEST_APP_SERVICE"
	appServiceEnv := os.Getenv(appServiceEnvVariable)
	if appServiceEnv == "" {
		t.Skipf("Skipping as %q is not specified", appServiceEnvVariable)
	}

	domainEnvVariable := "ARM_TEST_DOMAIN"
	domainEnv := os.Getenv(domainEnvVariable)
	if domainEnv == "" {
		t.Skipf("Skipping as %q is not specified", domainEnvVariable)
	}

	resourceName := "azurerm_app_service_custom_hostname_binding.test"
	ri := acctest.RandInt()
	location := testLocation()
	config := testAccAzureRMAppServiceCustomHostnameBinding_basicConfig(ri, location, appServiceEnv, domainEnv)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceCustomHostnameBindingDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceCustomHostnameBindingExists(resourceName),
				),
			},
		},
	})
}

func testCheckAzureRMAppServiceCustomHostnameBindingDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).appServicesClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_app_service_custom_hostname_binding" {
			continue
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		appServiceName := rs.Primary.Attributes["app_service_name"]
		hostname := rs.Primary.Attributes["hostname"]

		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.GetHostNameBinding(ctx, resourceGroup, appServiceName, hostname)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}
			return err
		}

		return nil
	}

	return nil
}

func testCheckAzureRMAppServiceCustomHostnameBindingExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		appServiceName := rs.Primary.Attributes["app_service_name"]
		hostname := rs.Primary.Attributes["hostname"]

		client := testAccProvider.Meta().(*ArmClient).appServicesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.GetHostNameBinding(ctx, resourceGroup, appServiceName, hostname)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Hostname Binding %q (App Service %q / Resource Group: %q) does not exist", hostname, appServiceName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on appServicesClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMAppServiceCustomHostnameBinding_basicConfig(rInt int, location string, appServiceName string, domain string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
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

resource "azurerm_app_service" "test" {
  name                = "%s"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
}

resource "azurerm_app_service_custom_hostname_binding" "test" {
  hostname            = "%s"
  app_service_name    = "${azurerm_app_service.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, rInt, location, rInt, appServiceName, domain)
}
