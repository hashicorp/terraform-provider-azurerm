package automation_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccAzureRMAutomationVariableBool_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_variable_bool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationVariableBoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationVariableBool_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationVariableBoolExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "value", "false"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAutomationVariableBool_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_variable_bool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationVariableBoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationVariableBool_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationVariableBoolExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "description", "This variable is created by Terraform acceptance test."),
					resource.TestCheckResourceAttr(data.ResourceName, "value", "true"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAutomationVariableBool_basicCompleteUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_variable_bool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationVariableBoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationVariableBool_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationVariableBoolExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "value", "false"),
				),
			},
			{
				Config: testAccAzureRMAutomationVariableBool_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationVariableBoolExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "description", "This variable is created by Terraform acceptance test."),
					resource.TestCheckResourceAttr(data.ResourceName, "value", "true"),
				),
			},
			{
				Config: testAccAzureRMAutomationVariableBool_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationVariableBoolExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "value", "false"),
				),
			},
		},
	})
}

func testCheckAzureRMAutomationVariableBoolExists(resourceName string) resource.TestCheckFunc {
	return testCheckAzureRMAutomationVariableExists(resourceName, "Bool")
}

func testCheckAzureRMAutomationVariableBoolDestroy(s *terraform.State) error {
	return testCheckAzureRMAutomationVariableDestroy(s, "Bool")
}

func testAccAzureRMAutomationVariableBool_basic(data acceptance.TestData) string {
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

resource "azurerm_automation_variable_bool" "test" {
  name                    = "acctestAutoVar-%d"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  value                   = false
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMAutomationVariableBool_complete(data acceptance.TestData) string {
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

resource "azurerm_automation_variable_bool" "test" {
  name                    = "acctestAutoVar-%d"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  description             = "This variable is created by Terraform acceptance test."
  value                   = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
