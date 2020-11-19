package automation_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccAzureRMAutomationVariableDateTime_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_variable_datetime", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationVariableDateTimeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationVariableDateTime_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationVariableDateTimeExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "value", "2019-04-24T21:40:54.074Z"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAutomationVariableDateTime_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_variable_datetime", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationVariableDateTimeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationVariableDateTime_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationVariableDateTimeExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "description", "This variable is created by Terraform acceptance test."),
					resource.TestCheckResourceAttr(data.ResourceName, "value", "2019-04-20T08:40:04.02Z"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAutomationVariableDateTime_basicCompleteUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_variable_datetime", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationVariableDateTimeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationVariableDateTime_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationVariableDateTimeExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "value", "2019-04-24T21:40:54.074Z"),
				),
			},
			{
				Config: testAccAzureRMAutomationVariableDateTime_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationVariableDateTimeExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "description", "This variable is created by Terraform acceptance test."),
					resource.TestCheckResourceAttr(data.ResourceName, "value", "2019-04-20T08:40:04.02Z"),
				),
			},
			{
				Config: testAccAzureRMAutomationVariableDateTime_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationVariableDateTimeExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "value", "2019-04-24T21:40:54.074Z"),
				),
			},
		},
	})
}

func testCheckAzureRMAutomationVariableDateTimeExists(resourceName string) resource.TestCheckFunc {
	return testCheckAzureRMAutomationVariableExists(resourceName, "Datetime")
}

func testCheckAzureRMAutomationVariableDateTimeDestroy(s *terraform.State) error {
	return testCheckAzureRMAutomationVariableDestroy(s, "Datetime")
}

func testAccAzureRMAutomationVariableDateTime_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctestAutoAcct-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}

resource "azurerm_automation_variable_datetime" "test" {
  name                    = "acctestAutoVar-%d"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  value                   = "2019-04-24T21:40:54.074Z"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMAutomationVariableDateTime_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctestAutoAcct-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}

resource "azurerm_automation_variable_datetime" "test" {
  name                    = "acctestAutoVar-%d"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  description             = "This variable is created by Terraform acceptance test."
  value                   = "2019-04-20T08:40:04.02Z"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
