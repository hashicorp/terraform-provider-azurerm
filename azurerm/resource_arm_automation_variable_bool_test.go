package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccAzureRMAutomationVariableBool_basic(t *testing.T) {
	resourceName := "azurerm_automation_variable_bool.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationVariableBoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationVariableBool_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationVariableBoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "value", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMAutomationVariableBool_complete(t *testing.T) {
	resourceName := "azurerm_automation_variable_bool.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationVariableBoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationVariableBool_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationVariableBoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", "This variable is created by Terraform acceptance test."),
					resource.TestCheckResourceAttr(resourceName, "value", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMAutomationVariableBool_basicCompleteUpdate(t *testing.T) {
	resourceName := "azurerm_automation_variable_bool.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationVariableBoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationVariableBool_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationVariableBoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "value", "false"),
				),
			},
			{
				Config: testAccAzureRMAutomationVariableBool_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationVariableBoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", "This variable is created by Terraform acceptance test."),
					resource.TestCheckResourceAttr(resourceName, "value", "true"),
				),
			},
			{
				Config: testAccAzureRMAutomationVariableBool_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationVariableBoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "value", "false"),
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

func testAccAzureRMAutomationVariableBool_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctestAutoAcct-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name = "Basic"
  }
}

resource "azurerm_automation_variable_bool" "test" {
  name                    = "acctestAutoVar-%d"
  resource_group_name     = "${azurerm_resource_group.test.name}"
  automation_account_name = "${azurerm_automation_account.test.name}"
  value                   = false
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAutomationVariableBool_complete(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctestAutoAcct-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name = "Basic"
  }
}

resource "azurerm_automation_variable_bool" "test" {
  name                    = "acctestAutoVar-%d"
  resource_group_name     = "${azurerm_resource_group.test.name}"
  automation_account_name = "${azurerm_automation_account.test.name}"
  description             = "This variable is created by Terraform acceptance test."
  value                   = true
}
`, rInt, location, rInt, rInt)
}
