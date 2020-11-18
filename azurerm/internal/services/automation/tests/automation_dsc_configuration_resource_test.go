package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMAutomationDscConfiguration_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_dsc_configuration", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationDscConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationDscConfiguration_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationDscConfigurationExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "location"),
					resource.TestCheckResourceAttr(data.ResourceName, "description", "test"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "log_verbose"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "state"),
					resource.TestCheckResourceAttr(data.ResourceName, "content_embedded", "configuration acctest {}"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.ENV", "prod"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAutomationDscConfiguration_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_dsc_configuration", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationDscConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationDscConfiguration_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationDscConfigurationExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMAutomationDscConfiguration_requiresImport),
		},
	})
}

func testCheckAzureRMAutomationDscConfigurationDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).Automation.DscConfigurationClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_automation_dsc_configuration" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		accName := rs.Primary.Attributes["automation_account_name"]

		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Automation Dsc Configuration: '%s'", name)
		}

		resp, err := conn.Get(ctx, resourceGroup, accName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("Automation Dsc Configuration still exists:\n%#v", resp)
	}

	return nil
}

func testCheckAzureRMAutomationDscConfigurationExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).Automation.DscConfigurationClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		accName := rs.Primary.Attributes["automation_account_name"]

		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Automation Dsc Configuration: '%s'", name)
		}

		resp, err := conn.Get(ctx, resourceGroup, accName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Automation Dsc Configuration '%s' (resource group: '%s') does not exist", name, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on automationDscConfigurationClient: %s\nName: %s, Account name: %s, Resource group: %s OBJECT: %+v", err, name, accName, resourceGroup, rs.Primary)
		}

		return nil
	}
}

func testAccAzureRMAutomationDscConfiguration_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctest-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}

resource "azurerm_automation_dsc_configuration" "test" {
  name                    = "acctest"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  location                = azurerm_resource_group.test.location
  content_embedded        = "configuration acctest {}"
  description             = "test"

  tags = {
    ENV = "prod"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMAutomationDscConfiguration_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMAutomationDscConfiguration_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_automation_dsc_configuration" "import" {
  name                    = azurerm_automation_dsc_configuration.test.name
  resource_group_name     = azurerm_automation_dsc_configuration.test.resource_group_name
  automation_account_name = azurerm_automation_dsc_configuration.test.automation_account_name
  location                = azurerm_automation_dsc_configuration.test.location
  content_embedded        = azurerm_automation_dsc_configuration.test.content_embedded
  description             = azurerm_automation_dsc_configuration.test.description
}
`, template)
}
