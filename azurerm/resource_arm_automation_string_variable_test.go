package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccAzureRMAutomationStringVariable_basic(t *testing.T) {
	resourceName := "azurerm_automation_string_variable.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutomationStringVariableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationStringVariable_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationStringVariableExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "value", "Hello, Terraform Basic Test."),
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

func TestAccAzureRMAutomationStringVariable_complete(t *testing.T) {
	resourceName := "azurerm_automation_string_variable.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutomationStringVariableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationStringVariable_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationStringVariableExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", "This variable is created by Terraform acceptance test."),
					resource.TestCheckResourceAttr(resourceName, "value", "Hello, Terraform Complete Test."),
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

func TestAccAzureRMAutomationStringVariable_basicCompleteUpdate(t *testing.T) {
	resourceName := "azurerm_automation_string_variable.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutomationStringVariableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationStringVariable_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationStringVariableExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "value", "Hello, Terraform Basic Test."),
				),
			},
			{
				Config: testAccAzureRMAutomationStringVariable_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationStringVariableExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", "This variable is created by Terraform acceptance test."),
					resource.TestCheckResourceAttr(resourceName, "value", "Hello, Terraform Complete Test."),
				),
			},
			{
				Config: testAccAzureRMAutomationStringVariable_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationStringVariableExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "value", "Hello, Terraform Basic Test."),
				),
			},
		},
	})
}

func testCheckAzureRMAutomationStringVariableExists(resourceName string) resource.TestCheckFunc {
	return testCheckAzureRMAutomationVariableExists(resourceName, "String")
}

func testCheckAzureRMAutomationStringVariableDestroy(s *terraform.State) error {
	return testCheckAzureRMAutomationVariableDestroy(s, "String")
}

func testAccAzureRMAutomationStringVariable_basic(rInt int, location string) string {
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

resource "azurerm_automation_string_variable" "test" {
  name                    = "acctestAutoVar-%d"
  resource_group_name     = "${azurerm_resource_group.test.name}"
  automation_account_name = "${azurerm_automation_account.test.name}"
  value                   = "Hello, Terraform Basic Test."
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAutomationStringVariable_complete(rInt int, location string) string {
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

resource "azurerm_automation_string_variable" "test" {
  name                    = "acctestAutoVar-%d"
  resource_group_name     = "${azurerm_resource_group.test.name}"
  automation_account_name = "${azurerm_automation_account.test.name}"
  description             = "This variable is created by Terraform acceptance test."
  value                   = "Hello, Terraform Complete Test."
}
`, rInt, location, rInt, rInt)
}
