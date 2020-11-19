package automation_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccAzureRMAutomationConnectionServicePrincipal_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_connection_service_principal", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationConnectionServicePrincipalDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationConnectionServicePrincipal_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationConnectionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAutomationConnectionServicePrincipal_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_connection_service_principal", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationConnectionServicePrincipal_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationConnectionExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMAutomationConnectionServicePrincipal_requiresImport),
		},
	})
}

func TestAccAzureRMAutomationConnectionServicePrincipal_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_connection_service_principal", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationConnectionServicePrincipalDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationConnectionServicePrincipal_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationConnectionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAutomationConnectionServicePrincipal_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_connection_service_principal", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationConnectionServicePrincipalDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationConnectionServicePrincipal_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationConnectionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMAutomationConnectionServicePrincipal_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationConnectionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMAutomationConnectionServicePrincipal_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationConnectionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMAutomationConnectionServicePrincipalDestroy(s *terraform.State) error {
	return testCheckAzureRMAutomationConnectionDestroy(s, "service_principal")
}

func testAccAzureRMAutomationConnectionServicePrincipal_basic(data acceptance.TestData) string {
	template := testAccAzureRMAutomationConnectionServicePrincipal_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_automation_connection_service_principal" "test" {
  name                    = "acctestACSP-%d"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  application_id          = "00000000-0000-0000-0000-000000000000"
  tenant_id               = data.azurerm_client_config.test.tenant_id
  subscription_id         = data.azurerm_client_config.test.subscription_id
  certificate_thumbprint  = file("testdata/automation_certificate_test.thumb")
}
`, template, data.RandomInteger)
}

func testAccAzureRMAutomationConnectionServicePrincipal_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMAutomationConnectionServicePrincipal_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_automation_connection_service_principal" "import" {
  name                    = azurerm_automation_connection_service_principal.test.name
  resource_group_name     = azurerm_automation_connection_service_principal.test.resource_group_name
  automation_account_name = azurerm_automation_connection_service_principal.test.automation_account_name
  application_id          = azurerm_automation_connection_service_principal.test.application_id
  tenant_id               = azurerm_automation_connection_service_principal.test.tenant_id
  subscription_id         = azurerm_automation_connection_service_principal.test.subscription_id
  certificate_thumbprint  = azurerm_automation_connection_service_principal.test.certificate_thumbprint
}
`, template)
}

func testAccAzureRMAutomationConnectionServicePrincipal_complete(data acceptance.TestData) string {
	template := testAccAzureRMAutomationConnectionServicePrincipal_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_automation_connection_service_principal" "test" {
  name                    = "acctestACSP-%d"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  application_id          = "00000000-0000-0000-0000-000000000000"
  tenant_id               = data.azurerm_client_config.test.tenant_id
  subscription_id         = data.azurerm_client_config.test.subscription_id
  certificate_thumbprint  = file("testdata/automation_certificate_test.thumb")
  description             = "acceptance test for automation connection"
}
`, template, data.RandomInteger)
}

func testAccAzureRMAutomationConnectionServicePrincipal_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-auto-%d"
  location = "%s"
}

data "azurerm_client_config" "test" {}

resource "azurerm_automation_account" "test" {
  name                = "acctestAA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
