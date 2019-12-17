package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMLogicAppActionCustom_basic(t *testing.T) {
	resourceName := "azurerm_logic_app_action_custom.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()
	config := testAccAzureRMLogicAppActionCustom_basic(ri, location)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppActionExists(resourceName),
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

func TestAccAzureRMLogicAppActionCustom_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_logic_app_action_custom.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
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
				ExpectError: acceptance.RequiresImportError("azurerm_logic_app_action_custom"),
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
  name         = "${azurerm_logic_app_action_custom.test.name}"
  logic_app_id = "${azurerm_logic_app_action_custom.test.logic_app_id}"
  body         = "${azurerm_logic_app_action_custom.test.body}"
}
`, template)
}

func testAccAzureRMLogicAppActionCustom_template(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_logic_app_workflow" "test" {
  name                = "acctestlaw-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, rInt, location, rInt)
}
