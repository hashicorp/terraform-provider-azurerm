package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccAzureRMAutomationVariableDateTime_basic(t *testing.T) {
	resourceName := "azurerm_automation_variable_datetime.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationVariableDateTimeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationVariableDateTime_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationVariableDateTimeExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "value", "2019-04-24T21:40:54.074Z"),
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

func TestAccAzureRMAutomationVariableDateTime_complete(t *testing.T) {
	resourceName := "azurerm_automation_variable_datetime.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationVariableDateTimeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationVariableDateTime_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationVariableDateTimeExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", "This variable is created by Terraform acceptance test."),
					resource.TestCheckResourceAttr(resourceName, "value", "2019-04-20T08:40:04.02Z"),
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

func TestAccAzureRMAutomationVariableDateTime_basicCompleteUpdate(t *testing.T) {
	resourceName := "azurerm_automation_variable_datetime.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationVariableDateTimeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationVariableDateTime_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationVariableDateTimeExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "value", "2019-04-24T21:40:54.074Z"),
				),
			},
			{
				Config: testAccAzureRMAutomationVariableDateTime_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationVariableDateTimeExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", "This variable is created by Terraform acceptance test."),
					resource.TestCheckResourceAttr(resourceName, "value", "2019-04-20T08:40:04.02Z"),
				),
			},
			{
				Config: testAccAzureRMAutomationVariableDateTime_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationVariableDateTimeExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "value", "2019-04-24T21:40:54.074Z"),
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

func testAccAzureRMAutomationVariableDateTime_basic(rInt int, location string) string {
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

resource "azurerm_automation_variable_datetime" "test" {
  name                    = "acctestAutoVar-%d"
  resource_group_name     = "${azurerm_resource_group.test.name}"
  automation_account_name = "${azurerm_automation_account.test.name}"
  value                   = "2019-04-24T21:40:54.074Z"
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAutomationVariableDateTime_complete(rInt int, location string) string {
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

resource "azurerm_automation_variable_datetime" "test" {
  name                    = "acctestAutoVar-%d"
  resource_group_name     = "${azurerm_resource_group.test.name}"
  automation_account_name = "${azurerm_automation_account.test.name}"
  description             = "This variable is created by Terraform acceptance test."
  value                   = "2019-04-20T08:40:04.02Z"
}
`, rInt, location, rInt, rInt)
}
