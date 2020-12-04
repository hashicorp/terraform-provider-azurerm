package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMSynapseWorkspace_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_synapse_workspace", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSynapseWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSynapseWorkspace_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSynapseWorkspaceExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "connectivity_endpoints.%"),
				),
			},
		},
	})
}

func testAccDataSourceSynapseWorkspace_basic(data acceptance.TestData) string {
	config := testAccAzureRMSynapseWorkspace_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_synapse_workspace" "test" {
  name                = azurerm_synapse_workspace.test.name
  resource_group_name = azurerm_synapse_workspace.test.resource_group_name
}
`, config)
}
