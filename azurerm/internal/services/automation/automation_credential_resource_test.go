package automation_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMAutomationCredential_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_credential", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationCredential_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationCredentialExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "username", "test_user"),
				),
			},
			data.ImportStep("password"),
		},
	})
}

func TestAccAzureRMAutomationCredential_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_credential", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationCredential_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationCredentialExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMAutomationCredential_requiresImport),
		},
	})
}

func TestAccAzureRMAutomationCredential_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_credential", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationCredential_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationCredentialExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "username", "test_user"),
					resource.TestCheckResourceAttr(data.ResourceName, "description", "This is a test credential for terraform acceptance test"),
				),
			},
			data.ImportStep("password"),
		},
	})
}

func testCheckAzureRMAutomationCredentialDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).Automation.CredentialClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_automation_credential" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		accName := rs.Primary.Attributes["automation_account_name"]

		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Automation Credential: '%s'", name)
		}

		resp, err := conn.Get(ctx, resourceGroup, accName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("Automation Credential still exists:\n%#v", resp)
	}

	return nil
}

func testCheckAzureRMAutomationCredentialExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).Automation.CredentialClient
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
			return fmt.Errorf("Bad: no resource group found in state for Automation Credential: '%s'", name)
		}

		resp, err := conn.Get(ctx, resourceGroup, accName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Automation Credential '%s' (resource group: '%s') does not exist", name, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on automationCredentialClient: %s\nName: %s, Account name: %s, Resource group: %s OBJECT: %+v", err, name, accName, resourceGroup, rs.Primary)
		}

		return nil
	}
}

func testAccAzureRMAutomationCredential_basic(data acceptance.TestData) string {
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

resource "azurerm_automation_credential" "test" {
  name                    = "acctest-%d"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  username                = "test_user"
  password                = "test_pwd"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMAutomationCredential_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMAutomationCredential_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_automation_credential" "import" {
  name                    = azurerm_automation_credential.test.name
  resource_group_name     = azurerm_automation_credential.test.resource_group_name
  automation_account_name = azurerm_automation_credential.test.automation_account_name
  username                = azurerm_automation_credential.test.username
  password                = azurerm_automation_credential.test.password
}
`, template)
}

func testAccAzureRMAutomationCredential_complete(data acceptance.TestData) string {
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

resource "azurerm_automation_credential" "test" {
  name                    = "acctest-%d"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  username                = "test_user"
  password                = "test_pwd"
  description             = "This is a test credential for terraform acceptance test"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
