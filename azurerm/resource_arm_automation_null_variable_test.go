package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccAzureRMAutomationNullVariable_basic(t *testing.T) {
	resourceName := "azurerm_automation_null_variable.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutomationNullVariableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationNullVariable_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationNullVariableExists(resourceName),
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

func TestAccAzureRMAutomationNullVariable_complete(t *testing.T) {
	resourceName := "azurerm_automation_null_variable.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutomationNullVariableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationNullVariable_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationNullVariableExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", "This variable is created by Terraform acceptance test."),
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

func TestAccAzureRMAutomationNullVariable_basicCompleteUpdate(t *testing.T) {
	resourceName := "azurerm_automation_null_variable.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutomationNullVariableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationNullVariable_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationNullVariableExists(resourceName),
				),
			},
			{
				Config: testAccAzureRMAutomationNullVariable_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationNullVariableExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", "This variable is created by Terraform acceptance test."),
				),
			},
			{
				Config: testAccAzureRMAutomationNullVariable_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationNullVariableExists(resourceName),
				),
			},
		},
	})
}

func testCheckAzureRMAutomationNullVariableExists(resourceName string) resource.TestCheckFunc {
	return testCheckAzureRMAutomationVariableExists(resourceName, "Null")
}

func testCheckAzureRMAutomationNullVariableDestroy(s *terraform.State) error {
	return testCheckAzureRMAutomationVariableDestroy(s, "Null")
}

func testAccAzureRMAutomationNullVariable_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctestAutoAcct-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku = {
    name = "Basic"
  }
}

resource "azurerm_automation_null_variable" "test" {
  name                    = "acctestAutoVar-%d"
  resource_group_name     = "${azurerm_resource_group.test.name}"
  automation_account_name = "${azurerm_automation_account.test.name}"
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAutomationNullVariable_complete(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctestAutoAcct-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku = {
    name = "Basic"
  }
}

resource "azurerm_automation_null_variable" "test" {
  name                    = "acctestAutoVar-%d"
  resource_group_name     = "${azurerm_resource_group.test.name}"
  automation_account_name = "${azurerm_automation_account.test.name}"
  description             = "This variable is created by Terraform acceptance test."
}
`, rInt, location, rInt, rInt)
}
