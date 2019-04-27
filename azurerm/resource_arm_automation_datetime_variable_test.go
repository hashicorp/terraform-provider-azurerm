package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
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
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Automation Datetime Variable not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		accountName := rs.Primary.Attributes["automation_account_name"]

		client := testAccProvider.Meta().(*ArmClient).automationVariableClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		if resp, err := client.Get(ctx, resourceGroup, accountName, name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Automation Datetime Variable %q (Automation Account Name %q / Resource Group %q) does not exist", name, accountName, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on automationVariableClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMAutomationDatetimeVariableDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).automationVariableClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_automation_datetime_variable" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		accountName := rs.Primary.Attributes["automation_account_name"]

		if resp, err := client.Get(ctx, resourceGroup, accountName, name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on automationVariableClient: %+v", err)
			}
		}

		return nil
	}

	return nil
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
