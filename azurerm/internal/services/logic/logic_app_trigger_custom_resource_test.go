package logic_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccAzureRMLogicAppTriggerCustom_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_trigger_custom", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogicAppTriggerCustom_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppTriggerExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLogicAppTriggerCustom_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_trigger_custom", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogicAppTriggerCustom_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppTriggerExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMLogicAppTriggerCustom_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_logic_app_trigger_custom"),
			},
		},
	})
}

func testAccAzureRMLogicAppTriggerCustom_basic(data acceptance.TestData) string {
	template := testAccAzureRMLogicAppTriggerCustom_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_trigger_custom" "test" {
  name         = "recurrence-%d"
  logic_app_id = azurerm_logic_app_workflow.test.id

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
`, template, data.RandomInteger)
}

func testAccAzureRMLogicAppTriggerCustom_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMLogicAppTriggerCustom_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_trigger_custom" "import" {
  name         = azurerm_logic_app_trigger_custom.test.name
  logic_app_id = azurerm_logic_app_trigger_custom.test.logic_app_id
  body         = azurerm_logic_app_trigger_custom.test.body
}
`, template)
}

func testAccAzureRMLogicAppTriggerCustom_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_logic_app_workflow" "test" {
  name                = "acctestlaw-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
