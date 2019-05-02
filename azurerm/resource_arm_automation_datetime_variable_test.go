package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccAzureRMAutomationDatetimeVariable_basic(t *testing.T) {
	resourceName := "azurerm_automation_datetime_variable.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutomationDatetimeVariableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationDatetimeVariable_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationDatetimeVariableExists(resourceName),
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

func TestAccAzureRMAutomationDatetimeVariable_complete(t *testing.T) {
	resourceName := "azurerm_automation_datetime_variable.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutomationDatetimeVariableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationDatetimeVariable_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationDatetimeVariableExists(resourceName),
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

func TestAccAzureRMAutomationDatetimeVariable_basicCompleteUpdate(t *testing.T) {
	resourceName := "azurerm_automation_datetime_variable.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutomationDatetimeVariableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationDatetimeVariable_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationDatetimeVariableExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "value", "2019-04-24T21:40:54.074Z"),
				),
			},
			{
				Config: testAccAzureRMAutomationDatetimeVariable_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationDatetimeVariableExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", "This variable is created by Terraform acceptance test."),
					resource.TestCheckResourceAttr(resourceName, "value", "2019-04-20T08:40:04.02Z"),
				),
			},
			{
				Config: testAccAzureRMAutomationDatetimeVariable_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationDatetimeVariableExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "value", "2019-04-24T21:40:54.074Z"),
				),
			},
		},
	})
}

func testCheckAzureRMAutomationDatetimeVariableExists(resourceName string) resource.TestCheckFunc {
	return testCheckAzureRMAutomationVariableExists(resourceName, "Datetime")
}

func testCheckAzureRMAutomationDatetimeVariableDestroy(s *terraform.State) error {
	return testCheckAzureRMAutomationVariableDestroy(s, "Datetime")
}

func testAccAzureRMAutomationDatetimeVariable_basic(rInt int, location string) string {
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

resource "azurerm_automation_datetime_variable" "test" {
  name                    = "acctestAutoVar-%d"
  resource_group_name     = "${azurerm_resource_group.test.name}"
  automation_account_name = "${azurerm_automation_account.test.name}"
  value                   = "2019-04-24T21:40:54.074Z"
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAutomationDatetimeVariable_complete(rInt int, location string) string {
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

resource "azurerm_automation_datetime_variable" "test" {
  name                    = "acctestAutoVar-%d"
  resource_group_name     = "${azurerm_resource_group.test.name}"
  automation_account_name = "${azurerm_automation_account.test.name}"
  description             = "This variable is created by Terraform acceptance test."
  value                   = "2019-04-20T08:40:04.02Z"
}
`, rInt, location, rInt, rInt)
}
