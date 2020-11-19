package automation_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccAzureRMAutomationConnectionClassicCertificate_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_connection_classic_certificate", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationConnectionClassicCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationConnectionClassicCertificate_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationConnectionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAutomationConnectionClassicCertificate_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_connection_classic_certificate", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationConnectionClassicCertificate_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationConnectionExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMAutomationConnectionClassicCertificate_requiresImport),
		},
	})
}

func TestAccAzureRMAutomationConnectionClassicCertificate_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_connection_classic_certificate", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationConnectionClassicCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationConnectionClassicCertificate_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationConnectionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAutomationConnectionClassicCertificate_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_connection_classic_certificate", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationConnectionClassicCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationConnectionClassicCertificate_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationConnectionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMAutomationConnectionClassicCertificate_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationConnectionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMAutomationConnectionClassicCertificate_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationConnectionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMAutomationConnectionClassicCertificateDestroy(s *terraform.State) error {
	return testCheckAzureRMAutomationConnectionDestroy(s, "classic_certificate")
}

func testAccAzureRMAutomationConnectionClassicCertificate_basic(data acceptance.TestData) string {
	template := testAccAzureRMAutomationConnectionClassicCertificate_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_automation_connection_classic_certificate" "test" {
  name                    = "acctestACCC-%d"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  certificate_asset_name  = "cert1"
  subscription_name       = "subs1"
  subscription_id         = data.azurerm_client_config.test.subscription_id
}
`, template, data.RandomInteger)
}

func testAccAzureRMAutomationConnectionClassicCertificate_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMAutomationConnectionClassicCertificate_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_automation_connection_classic_certificate" "import" {
  name                    = azurerm_automation_connection_classic_certificate.test.name
  resource_group_name     = azurerm_automation_connection_classic_certificate.test.resource_group_name
  automation_account_name = azurerm_automation_connection_classic_certificate.test.automation_account_name
  certificate_asset_name  = azurerm_automation_connection_classic_certificate.test.certificate_asset_name
  subscription_name       = azurerm_automation_connection_classic_certificate.test.subscription_name
  subscription_id         = azurerm_automation_connection_classic_certificate.test.subscription_id
}
`, template)
}

func testAccAzureRMAutomationConnectionClassicCertificate_complete(data acceptance.TestData) string {
	template := testAccAzureRMAutomationConnectionClassicCertificate_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_automation_connection_classic_certificate" "test" {
  name                    = "acctestACCC-%d"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  certificate_asset_name  = "cert1"
  subscription_name       = "subs1"
  subscription_id         = data.azurerm_client_config.test.subscription_id
  description             = "acceptance test for automation connection"
}
`, template, data.RandomInteger)
}

func testAccAzureRMAutomationConnectionClassicCertificate_template(data acceptance.TestData) string {
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
