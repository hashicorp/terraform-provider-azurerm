package logic_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccAzureRMLogicAppActionCustom_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_action_custom", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogicAppActionCustom_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppActionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLogicAppActionCustom_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_action_custom", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogicAppActionCustom_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppActionExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMLogicAppActionCustom_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_logic_app_action_custom"),
			},
		},
	})
}

func testAccAzureRMLogicAppActionCustom_basic(data acceptance.TestData) string {
	template := testAccAzureRMLogicAppActionCustom_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_action_custom" "test" {
  name         = "action%d"
  logic_app_id = azurerm_logic_app_workflow.test.id

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
`, template, data.RandomInteger)
}

func testAccAzureRMLogicAppActionCustom_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMLogicAppActionCustom_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_action_custom" "import" {
  name         = azurerm_logic_app_action_custom.test.name
  logic_app_id = azurerm_logic_app_action_custom.test.logic_app_id
  body         = azurerm_logic_app_action_custom.test.body
}
`, template)
}

func testAccAzureRMLogicAppActionCustom_template(data acceptance.TestData) string {
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
