package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMLogicAppWorkflow_basic(t *testing.T) {
	dataSourceName := "data.azurerm_logic_app_workflow.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()
	config := testAccDataSourceAzureRMLogicAppWorkflow_basic(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppWorkflowExists(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "parameters.%", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMLogicAppWorkflow_tags(t *testing.T) {
	dataSourceName := "data.azurerm_logic_app_workflow.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()
	config := testAccDataSourceAzureRMLogicAppWorkflow_tags(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppWorkflowExists(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "parameters.%", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.Source", "AcceptanceTests"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMLogicAppWorkflow_basic(rInt int, location string) string {
	r := testAccAzureRMLogicAppWorkflow_empty(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_logic_app_workflow" "test" {
  name                = "${azurerm_logic_app_workflow.test.name}"
  resource_group_name = "${azurerm_logic_app_workflow.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMLogicAppWorkflow_tags(rInt int, location string) string {
	r := testAccAzureRMLogicAppWorkflow_tags(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_logic_app_workflow" "test" {
  name                = "${azurerm_logic_app_workflow.test.name}"
  resource_group_name = "${azurerm_logic_app_workflow.test.resource_group_name}"
}
`, r)
}
