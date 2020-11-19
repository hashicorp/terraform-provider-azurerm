package automation_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccAzureRMAutomationVariableString_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_variable_string", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationVariableStringDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationVariableString_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationVariableStringExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "value", "Hello, Terraform Basic Test."),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAutomationVariableString_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_variable_string", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationVariableStringDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationVariableString_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationVariableStringExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "description", "This variable is created by Terraform acceptance test."),
					resource.TestCheckResourceAttr(data.ResourceName, "value", "Hello, Terraform Complete Test."),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAutomationVariableString_basicCompleteUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_variable_string", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationVariableStringDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationVariableString_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationVariableStringExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "value", "Hello, Terraform Basic Test."),
				),
			},
			{
				Config: testAccAzureRMAutomationVariableString_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationVariableStringExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "description", "This variable is created by Terraform acceptance test."),
					resource.TestCheckResourceAttr(data.ResourceName, "value", "Hello, Terraform Complete Test."),
				),
			},
			{
				Config: testAccAzureRMAutomationVariableString_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationVariableStringExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "value", "Hello, Terraform Basic Test."),
				),
			},
		},
	})
}

func testCheckAzureRMAutomationVariableStringExists(resourceName string) resource.TestCheckFunc {
	return testCheckAzureRMAutomationVariableExists(resourceName, "String")
}

func testCheckAzureRMAutomationVariableStringDestroy(s *terraform.State) error {
	return testCheckAzureRMAutomationVariableDestroy(s, "String")
}

func testAccAzureRMAutomationVariableString_basic(data acceptance.TestData) string {
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

resource "azurerm_automation_variable_string" "test" {
  name                    = "acctestAutoVar-%d"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  value                   = "Hello, Terraform Basic Test."
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMAutomationVariableString_complete(data acceptance.TestData) string {
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

resource "azurerm_automation_variable_string" "test" {
  name                    = "acctestAutoVar-%d"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  description             = "This variable is created by Terraform acceptance test."
  value                   = "Hello, Terraform Complete Test."
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
