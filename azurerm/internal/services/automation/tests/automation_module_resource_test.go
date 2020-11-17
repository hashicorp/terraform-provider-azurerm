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

func TestAccAzureRMAutomationModule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_module", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationModuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationModule_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationModuleExists(data.ResourceName),
				),
			},
			data.ImportStep("module_link"),
		},
	})
}

func TestAccAzureRMAutomationModule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_module", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationModuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationModule_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationModuleExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMAutomationModule_requiresImport),
		},
	})
}

func TestAccAzureRMAutomationModule_multipleModules(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_module", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationModuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationModule_multipleModules(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationModuleExists(data.ResourceName),
				),
			},
			data.ImportStep("module_link"),
		},
	})
}

func testCheckAzureRMAutomationModuleDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).Automation.ModuleClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_automation_module" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		accName := rs.Primary.Attributes["automation_account_name"]

		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Automation Module: '%s'", name)
		}

		resp, err := conn.Get(ctx, resourceGroup, accName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("Automation Module still exists:\n%#v", resp)
	}

	return nil
}

func testCheckAzureRMAutomationModuleExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).Automation.ModuleClient
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
			return fmt.Errorf("Bad: no resource group found in state for Automation Module: '%s'", name)
		}

		resp, err := conn.Get(ctx, resourceGroup, accName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Automation Module '%s' (resource group: '%s') does not exist", name, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on automationModuleClient: %s\nName: %s, Account name: %s, Resource group: %s OBJECT: %+v", err, name, accName, resourceGroup, rs.Primary)
		}

		return nil
	}
}

func testAccAzureRMAutomationModule_basic(data acceptance.TestData) string {
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

resource "azurerm_automation_module" "test" {
  name                    = "xActiveDirectory"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name

  module_link {
    uri = "https://devopsgallerystorage.blob.core.windows.net/packages/xactivedirectory.2.19.0.nupkg"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMAutomationModule_multipleModules(data acceptance.TestData) string {
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

resource "azurerm_automation_module" "test" {
  name                    = "xActiveDirectory"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name

  module_link {
    uri = "https://devopsgallerystorage.blob.core.windows.net/packages/xactivedirectory.2.19.0.nupkg"
  }
}

resource "azurerm_automation_module" "second" {
  name                    = "AzureRmMinus"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name

  module_link {
    uri = "https://www.powershellgallery.com/api/v2/package/AzureRmMinus/0.3.0.0"
  }

  depends_on = [azurerm_automation_module.test]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMAutomationModule_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMAutomationModule_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_automation_module" "import" {
  name                    = azurerm_automation_module.test.name
  resource_group_name     = azurerm_automation_module.test.resource_group_name
  automation_account_name = azurerm_automation_module.test.automation_account_name

  module_link {
    uri = "https://devopsgallerystorage.blob.core.windows.net/packages/xactivedirectory.2.19.0.nupkg"
  }
}
`, template)
}
