package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMLogicAppActionCustom_basic(t *testing.T) {
	resourceName := "azurerm_logic_app_action_custom.test"
	ri := acctest.RandInt()
	config := testAccAzureRMLogicAppActionCustom_basic(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppActionExists(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMLogicAppActionCustom_requiresImport(t *testing.T) {
	resourceName := "azurerm_logic_app_action_custom.test"
	ri := acctest.RandInt()
	location := testLocation()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogicAppActionCustom_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppActionExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMLogicAppActionCustom_requiresImport(ri, location),
				ExpectError: testRequiresImportError("azurerm_logic_app_action_custom"),
			},
		},
	})
}

func testAccAzureRMLogicAppActionCustom_basic(rInt int, location string) string {
	template := testAccAzureRMLogicAppActionCustom_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_action_custom" "test" {
  name         = "action%d"
  logic_app_id = "${azurerm_logic_app_workflow.test.id}"
  body = <<BODY
{
    "description": "A variable to configure the auto expiration age in days. Configured in negative number. Default is -30 (30 days old).",
    "inputs": {
        "variables": [
            {
                "name": "ExpirationAgeInDays",
                "type": "Integer",
                "value": -30
            }
        ]
    },
    "runAfter": {},
    "type": "InitializeVariable"
}
BODY
}
`, template, rInt)
}

func testAccAzureRMLogicAppActionCustom_requiresImport(rInt int, location string) string {
	template := testAccAzureRMLogicAppActionCustom_basic(rInt, location)

	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_action_custom" "import" {
  name         = "${azurerm_logic_app_action_custom.test.id}"
  logic_app_id = "${azurerm_logic_app_action_custom.test.logic_app_id}"
  body         = "${azurerm_logic_app_action_custom.test.body}"
}
`, template)
}

func testAccAzureRMLogicAppActionCustom_template(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_logic_app_workflow" "test" {
  name = "acctestlaw-%d"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, rInt, location, rInt)
}
