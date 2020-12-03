package logic_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMLogicAppWorkflow_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_logic_app_workflow", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMLogicAppWorkflow_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppWorkflowExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "parameters.%", "0"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "connector_endpoint_ip_addresses.#"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "connector_outbound_ip_addresses.#"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "workflow_endpoint_ip_addresses.#"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "workflow_outbound_ip_addresses.#"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMLogicAppWorkflow_tags(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_logic_app_workflow", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMLogicAppWorkflow_tags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppWorkflowExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "parameters.%", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.Source", "AcceptanceTests"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMLogicAppWorkflow_basic(data acceptance.TestData) string {
	r := testAccAzureRMLogicAppWorkflow_empty(data)
	return fmt.Sprintf(`
%s

data "azurerm_logic_app_workflow" "test" {
  name                = azurerm_logic_app_workflow.test.name
  resource_group_name = azurerm_logic_app_workflow.test.resource_group_name
}
`, r)
}

func testAccDataSourceAzureRMLogicAppWorkflow_tags(data acceptance.TestData) string {
	r := testAccAzureRMLogicAppWorkflow_tags(data)
	return fmt.Sprintf(`
%s

data "azurerm_logic_app_workflow" "test" {
  name                = azurerm_logic_app_workflow.test.name
  resource_group_name = azurerm_logic_app_workflow.test.resource_group_name
}
`, r)
}
