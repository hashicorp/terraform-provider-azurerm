package logic_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

type LogicAppTriggerCustomResource struct {
}

func TestAccLogicAppTriggerCustom_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_trigger_custom", "test")
	r := LogicAppTriggerCustomResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogicAppTriggerCustom_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_trigger_custom", "test")
	r := LogicAppTriggerCustomResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_logic_app_trigger_custom"),
		},
	})
}

func (LogicAppTriggerCustomResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	return triggerExists(ctx, clients, state)
}

func (LogicAppTriggerCustomResource) basic(data acceptance.TestData) string {
	template := LogicAppTriggerCustomResource{}.template(data)
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

func (LogicAppTriggerCustomResource) requiresImport(data acceptance.TestData) string {
	template := LogicAppTriggerCustomResource{}.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_trigger_custom" "import" {
  name         = azurerm_logic_app_trigger_custom.test.name
  logic_app_id = azurerm_logic_app_trigger_custom.test.logic_app_id
  body         = azurerm_logic_app_trigger_custom.test.body
}
`, template)
}

func (LogicAppTriggerCustomResource) template(data acceptance.TestData) string {
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
