package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMResourceGraphGraphQuery_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_resource_graph_graph_query", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMResourceGraphGraphQueryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceResourceGraphGraphQuery_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMResourceGraphGraphQueryExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "resource_name"),
				),
			},
		},
	})
}

func testAccDataSourceResourceGraphGraphQuery_basic(data acceptance.TestData) string {
	config := testAccAzureRMResourceGraphGraphQuery_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_resource_graph_graph_query" "test" {
  resource_group_name = azurerm_resource_graph_graph_query.test.resource_group_name
  resource_name       = azurerm_resource_graph_graph_query.test.resource_name
  query               = azurerm_resource_graph_graph_query.test.query
}
`, config)
}
