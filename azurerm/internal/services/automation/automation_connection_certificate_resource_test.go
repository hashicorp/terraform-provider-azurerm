package automation_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccAzureRMAutomationConnectionCertificate_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_connection_certificate", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationConnectionCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationConnectionCertificate_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationConnectionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAutomationConnectionCertificate_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_connection_certificate", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationConnectionCertificate_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationConnectionExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMAutomationConnectionCertificate_requiresImport),
		},
	})
}

func TestAccAzureRMAutomationConnectionCertificate_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_connection_certificate", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationConnectionCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationConnectionCertificate_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationConnectionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAutomationConnectionCertificate_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_connection_certificate", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationConnectionCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationConnectionCertificate_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationConnectionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMAutomationConnectionCertificate_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationConnectionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMAutomationConnectionCertificate_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationConnectionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMAutomationConnectionCertificateDestroy(s *terraform.State) error {
	return testCheckAzureRMAutomationConnectionDestroy(s, "certificate")
}

func testAccAzureRMAutomationConnectionCertificate_basic(data acceptance.TestData) string {
	template := testAccAzureRMAutomationConnectionCertificate_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_automation_connection_certificate" "test" {
  name                        = "acctestACC-%d"
  resource_group_name         = azurerm_resource_group.test.name
  automation_account_name     = azurerm_automation_account.test.name
  automation_certificate_name = azurerm_automation_certificate.test.name
  subscription_id             = data.azurerm_client_config.test.subscription_id
}
`, template, data.RandomInteger)
}

func testAccAzureRMAutomationConnectionCertificate_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMAutomationConnectionCertificate_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_automation_connection_certificate" "import" {
  name                        = azurerm_automation_connection_certificate.test.name
  resource_group_name         = azurerm_automation_connection_certificate.test.resource_group_name
  automation_account_name     = azurerm_automation_connection_certificate.test.automation_account_name
  automation_certificate_name = azurerm_automation_connection_certificate.test.automation_certificate_name
  subscription_id             = azurerm_automation_connection_certificate.test.subscription_id
}
`, template)
}

func testAccAzureRMAutomationConnectionCertificate_complete(data acceptance.TestData) string {
	template := testAccAzureRMAutomationConnectionCertificate_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_automation_connection_certificate" "test" {
  name                        = "acctestACC-%d"
  resource_group_name         = azurerm_resource_group.test.name
  automation_account_name     = azurerm_automation_account.test.name
  automation_certificate_name = azurerm_automation_certificate.test.name
  subscription_id             = data.azurerm_client_config.test.subscription_id
  description                 = "acceptance test for automation connection"
}
`, template, data.RandomInteger)
}

func testAccAzureRMAutomationConnectionCertificate_template(data acceptance.TestData) string {
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

resource "azurerm_automation_certificate" "test" {
  name                    = "acctest-%d"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  base64                  = filebase64("testdata/automation_certificate_test.pfx")
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
