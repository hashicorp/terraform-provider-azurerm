package azurerm

import (
	"fmt"
	"testing"

	"os"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMAppServiceCustomHostnameBinding(t *testing.T) {
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

	// NOTE: this is a combined test rather than separate split out tests due to
	// the app service name being shared (so the tests don't conflict with each other)
	testCases := map[string]map[string]func(t *testing.T, appServiceEnv, domainEnv string){
		"basic": {
			"basic":          testAccAzureRMAppServiceCustomHostnameBinding_basic,
			"multiple":       testAccAzureRMAppServiceCustomHostnameBinding_multiple,
			"requiresImport": testAccAzureRMAppServiceCustomHostnameBinding_requiresImport,
		},
	}

	for group, m := range testCases {
		m := m
		t.Run(group, func(t *testing.T) {
			for name, tc := range m {
				tc := tc
				t.Run(name, func(t *testing.T) {
					tc(t, appServiceEnv, domainEnv)
				})
			}
		})
	}
}

func testAccAzureRMAppServiceCustomHostnameBinding_basic(t *testing.T, appServiceEnv, domainEnv string) {
	resourceName := "azurerm_app_service_custom_hostname_binding.test"
	ri := tf.AccRandTimeInt()
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
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAzureRMAppServiceCustomHostnameBinding_requiresImport(t *testing.T, appServiceEnv, domainEnv string) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_app_service_custom_hostname_binding.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceCustomHostnameBindingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceCustomHostnameBinding_basicConfig(ri, location, appServiceEnv, domainEnv),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceCustomHostnameBindingExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMAppServiceCustomHostnameBinding_requiresImportConfig(ri, location, appServiceEnv, domainEnv),
				ExpectError: testRequiresImportError("azurerm_app_service_custom_hostname_binding"),
			},
		},
	})
}

func testAccAzureRMAppServiceCustomHostnameBinding_multiple(t *testing.T, appServiceEnv, domainEnv string) {
	altDomainEnvVariable := "ARM_ALT_TEST_DOMAIN"
	altDomainEnv := os.Getenv(altDomainEnvVariable)
	if altDomainEnv == "" {
		t.Skipf("Skipping as %q is not specified", altDomainEnvVariable)
	}

	resourceName := "azurerm_app_service_custom_hostname_binding.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	config := testAccAzureRMAppServiceCustomHostnameBinding_multipleConfig(ri, location, appServiceEnv, domainEnv, altDomainEnv)

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
	client := testAccProvider.Meta().(*ArmClient).web.AppServicesClient

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

func testCheckAzureRMAppServiceCustomHostnameBindingExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		appServiceName := rs.Primary.Attributes["app_service_name"]
		hostname := rs.Primary.Attributes["hostname"]

		client := testAccProvider.Meta().(*ArmClient).web.AppServicesClient
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
  name     = "acctestRG-%d"
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

func testAccAzureRMAppServiceCustomHostnameBinding_requiresImportConfig(rInt int, location string, appServiceName string, domain string) string {
	template := testAccAzureRMAppServiceCustomHostnameBinding_basicConfig(rInt, location, appServiceName, domain)
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_custom_hostname_binding" "import" {
  hostname            = "${azurerm_app_service_custom_hostname_binding.test.name}"
  app_service_name    = "${azurerm_app_service_custom_hostname_binding.test.app_service_name}"
  resource_group_name = "${azurerm_app_service_custom_hostname_binding.test.resource_group_name}"
}
`, template)
}

func testAccAzureRMAppServiceCustomHostnameBinding_multipleConfig(rInt int, location, appServiceName, domain, altDomain string) string {
	template := testAccAzureRMAppServiceCustomHostnameBinding_basicConfig(rInt, location, appServiceName, domain)
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_custom_hostname_binding" "test2" {
  hostname            = "%s"
  app_service_name    = "${azurerm_app_service.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, template, altDomain)
}
