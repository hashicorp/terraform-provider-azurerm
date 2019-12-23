package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMLogicAppTriggerCustom_basic(t *testing.T) {
	resourceName := "azurerm_logic_app_trigger_custom.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogicAppTriggerCustom_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppTriggerExists(resourceName),
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

func TestAccAzureRMLogicAppTriggerCustom_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_logic_app_trigger_custom.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogicAppTriggerCustom_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppTriggerExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMLogicAppTriggerCustom_requiresImport(ri, location),
				ExpectError: acceptance.RequiresImportError("azurerm_logic_app_trigger_custom"),
			},
		},
	})
}

func testAccAzureRMLogicAppTriggerCustom_basic(rInt int, location string) string {
	template := testAccAzureRMLogicAppTriggerCustom_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_trigger_custom" "test" {
  name         = "recurrence-%d"
  logic_app_id = "${azurerm_logic_app_workflow.test.id}"

  body = <<BODY
{
  "recurrence": {
    "frequency": "Day",
    "interval": 1
  },
  "type": "Recurrence"
}
BODY
}
`, template, rInt)
}

func testAccAzureRMLogicAppTriggerCustom_requiresImport(rInt int, location string) string {
	template := testAccAzureRMLogicAppTriggerCustom_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_trigger_custom" "import" {
  name         = "${azurerm_logic_app_trigger_custom.test.name}"
  logic_app_id = "${azurerm_logic_app_trigger_custom.test.logic_app_id}"
  body         = "${azurerm_logic_app_trigger_custom.test.body}"
}
`, template)
}

func testAccAzureRMLogicAppTriggerCustom_template(rInt int, location string) string {
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
