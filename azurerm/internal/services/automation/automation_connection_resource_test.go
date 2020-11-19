package automation_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccAzureRMAutomationConnection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_connection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationConnectionCustomDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationConnection_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationConnectionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAutomationConnection_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_connection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationConnectionCustomDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationConnection_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationConnectionExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMAutomationConnection_requiresImport),
		},
	})
}

func TestAccAzureRMAutomationConnection_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_connection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationConnectionCustomDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationConnection_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationConnectionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAutomationConnection_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_connection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationConnectionCustomDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationConnection_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationConnectionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMAutomationConnection_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationConnectionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMAutomationConnection_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationConnectionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMAutomationConnectionCustomDestroy(s *terraform.State) error {
	return testCheckAzureRMAutomationConnectionDestroy(s, "")
}

func testAccAzureRMAutomationConnection_basic(data acceptance.TestData) string {
	template := testAccAzureRMAutomationConnection_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_automation_connection" "test" {
  name                    = "acctestAAC-%d"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  type                    = "AzureServicePrincipal"

  values = {
    "ApplicationId" : "00000000-0000-0000-0000-000000000000"
    "TenantId" : data.azurerm_client_config.test.tenant_id
    "SubscriptionId" : data.azurerm_client_config.test.subscription_id
    "CertificateThumbprint" : file("testdata/automation_certificate_test.thumb")
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMAutomationConnection_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMAutomationConnection_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_automation_connection" "import" {
  name                    = azurerm_automation_connection.test.name
  resource_group_name     = azurerm_automation_connection.test.resource_group_name
  automation_account_name = azurerm_automation_connection.test.automation_account_name
  type                    = azurerm_automation_connection.test.type
  values                  = azurerm_automation_connection.test.values
}
`, template)
}

func testAccAzureRMAutomationConnection_complete(data acceptance.TestData) string {
	template := testAccAzureRMAutomationConnection_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_automation_connection" "test" {
  name                    = "acctestAAC-%d"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  type                    = "AzureServicePrincipal"
  description             = "acceptance test for automation connection"

  values = {
    "ApplicationId" : "00000000-0000-0000-0000-000000000000"
    "TenantId" : data.azurerm_client_config.test.tenant_id
    "SubscriptionId" : data.azurerm_client_config.test.subscription_id
    "CertificateThumbprint" : file("testdata/automation_certificate_test.thumb")
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMAutomationConnection_template(data acceptance.TestData) string {
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
