package kusto_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMKustoCluster_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kusto_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKustoClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMKustoCluster_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKustoClusterExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "uri"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "data_ingestion_uri"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMKustoCluster_basic(data acceptance.TestData) string {
	template := testAccAzureRMKustoCluster_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_kusto_cluster" "test" {
  name                = azurerm_kusto_cluster.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, template)
}
