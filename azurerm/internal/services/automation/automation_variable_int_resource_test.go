package automation_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccAzureRMAutomationVariableInt_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_variable_int", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationVariableIntDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationVariableInt_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationVariableIntExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "value", "1234"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAutomationVariableInt_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_variable_int", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationVariableIntDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationVariableInt_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationVariableIntExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "description", "This variable is created by Terraform acceptance test."),
					resource.TestCheckResourceAttr(data.ResourceName, "value", "12345"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAutomationVariableInt_basicCompleteUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_variable_int", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationVariableIntDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationVariableInt_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationVariableIntExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "value", "1234"),
				),
			},
			{
				Config: testAccAzureRMAutomationVariableInt_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationVariableIntExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "description", "This variable is created by Terraform acceptance test."),
					resource.TestCheckResourceAttr(data.ResourceName, "value", "12345"),
				),
			},
			{
				Config: testAccAzureRMAutomationVariableInt_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationVariableIntExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "value", "1234"),
				),
			},
		},
	})
}

func testCheckAzureRMAutomationVariableIntExists(resourceName string) resource.TestCheckFunc {
	return testCheckAzureRMAutomationVariableExists(resourceName, "Int")
}

func testCheckAzureRMAutomationVariableIntDestroy(s *terraform.State) error {
	return testCheckAzureRMAutomationVariableDestroy(s, "Int")
}

func testAccAzureRMAutomationVariableInt_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-auto-%d"
  location = "%s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctestAutoAcct-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}

resource "azurerm_automation_variable_int" "test" {
  name                    = "acctestAutoVar-%d"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  value                   = 1234
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMAutomationVariableInt_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-auto-%d"
  location = "%s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctestAutoAcct-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}

resource "azurerm_automation_variable_int" "test" {
  name                    = "acctestAutoVar-%d"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  description             = "This variable is created by Terraform acceptance test."
  value                   = 12345
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
