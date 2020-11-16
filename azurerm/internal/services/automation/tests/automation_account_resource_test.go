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

func TestAccAzureRMAutomationAccount_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationAccount_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku_name", "Basic"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "dsc_server_endpoint"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "dsc_primary_access_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "dsc_secondary_access_key"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAutomationAccount_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationAccount_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationAccountExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMAutomationAccount_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_automation_account"),
			},
		},
	})
}

func TestAccAzureRMAutomationAccount_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationAccount_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku_name", "Basic"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.hello", "world"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMAutomationAccountDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).Automation.AccountClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_automation_account" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("Automation Account still exists:\n%#v", resp)
	}

	return nil
}

func testCheckAzureRMAutomationAccountExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).Automation.AccountClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Automation Account: '%s'", name)
		}

		resp, err := conn.Get(ctx, resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Automation Account '%s' (resource group: '%s') was not found: %+v", name, resourceGroup, err)
			}

			return fmt.Errorf("Bad: Get on automationClient: %s", err)
		}

		return nil
	}
}

func testAccAzureRMAutomationAccount_basic(data acceptance.TestData) string {
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

  sku_name = "Basic"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMAutomationAccount_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMAutomationAccount_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_automation_account" "import" {
  name                = azurerm_automation_account.test.name
  location            = azurerm_automation_account.test.location
  resource_group_name = azurerm_automation_account.test.resource_group_name

  sku_name = "Basic"
}
`, template)
}

func testAccAzureRMAutomationAccount_complete(data acceptance.TestData) string {
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

  sku_name = "Basic"

  tags = {
    "hello" = "world"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
