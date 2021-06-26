package logic_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type LogicAppWorkflowDataSource struct {
}

func TestAccLogicAppWorkflowDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_logic_app_workflow", "test")
	r := LogicAppWorkflowDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("parameters.%").HasValue("0"),
				check.That(data.ResourceName).Key("connector_endpoint_ip_addresses.#").Exists(),
				check.That(data.ResourceName).Key("connector_outbound_ip_addresses.#").Exists(),
				check.That(data.ResourceName).Key("workflow_endpoint_ip_addresses.#").Exists(),
				check.That(data.ResourceName).Key("workflow_outbound_ip_addresses.#").Exists(),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
	})
}

func TestAccLogicAppWorkflowDataSource_tags(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_logic_app_workflow", "test")
	r := LogicAppWorkflowDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.tags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("parameters.%").HasValue("0"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.Source").HasValue("AcceptanceTests"),
			),
		},
	})
}

func (LogicAppWorkflowDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_logic_app_workflow" "test" {
  name                = azurerm_logic_app_workflow.test.name
  resource_group_name = azurerm_logic_app_workflow.test.resource_group_name
}
`, LogicAppWorkflowResource{}.empty(data))
}

func (LogicAppWorkflowDataSource) tags(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_logic_app_workflow" "test" {
  name                = azurerm_logic_app_workflow.test.name
  resource_group_name = azurerm_logic_app_workflow.test.resource_group_name
}
`, LogicAppWorkflowResource{}.tags(data))
}
