package logic_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/logic/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccLogicAppIntegrationAccount_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_integration_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogicAppIntegrationAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogicAppIntegrationAccount_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppIntegrationAccountExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccLogicAppIntegrationAccount_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_integration_account", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogicAppIntegrationAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogicAppIntegrationAccount_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppIntegrationAccountExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMLogicAppIntegrationAccount_requiresImport),
		},
	})
}

func TestAccLogicAppIntegrationAccount_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_integration_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogicAppIntegrationAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogicAppIntegrationAccount_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppIntegrationAccountExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccLogicAppIntegrationAccount_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_integration_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogicAppIntegrationAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogicAppIntegrationAccount_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppIntegrationAccountExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMLogicAppIntegrationAccount_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppIntegrationAccountExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMLogicAppIntegrationAccount_update(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppIntegrationAccountExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMLogicAppIntegrationAccount_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppIntegrationAccountExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMLogicAppIntegrationAccountExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Logic.IntegrationAccountClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Integration Account not found: %s", resourceName)
		}
		id, err := parse.IntegrationAccountID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.ResourceGroup, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: integration account %q does not exist", id.Name)
			}
			return fmt.Errorf("bad: Get on LogicAppIntegrationAccountClient: %+v", err)
		}
		return nil
	}
}

func testCheckAzureRMLogicAppIntegrationAccountDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Logic.IntegrationAccountClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_logic_app_integration_account" {
			continue
		}
		id, err := parse.IntegrationAccountID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.ResourceGroup, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Get on integration.accountClient: %+v", err)
			}
		}
		return nil
	}
	return nil
}

func testAccAzureRMLogicAppIntegrationAccount_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-logic-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMLogicAppIntegrationAccount_basic(data acceptance.TestData) string {
	template := testAccAzureRMLogicAppIntegrationAccount_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_integration_account" "test" {
  name                = "acctest-IA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}
`, template, data.RandomInteger)
}

func testAccAzureRMLogicAppIntegrationAccount_requiresImport(data acceptance.TestData) string {
	config := testAccAzureRMLogicAppIntegrationAccount_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_integration_account" "import" {
  name                = azurerm_logic_app_integration_account.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = azurerm_logic_app_integration_account.test.sku_name
}
`, config)
}

func testAccAzureRMLogicAppIntegrationAccount_complete(data acceptance.TestData) string {
	template := testAccAzureRMLogicAppIntegrationAccount_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_integration_account" "test" {
  name                = "acctest-IA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMLogicAppIntegrationAccount_update(data acceptance.TestData) string {
	template := testAccAzureRMLogicAppIntegrationAccount_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_integration_account" "test" {
  name                = "acctest-IA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Standard"
  tags = {
    ENV = "Stage"
  }
}
`, template, data.RandomInteger)
}
