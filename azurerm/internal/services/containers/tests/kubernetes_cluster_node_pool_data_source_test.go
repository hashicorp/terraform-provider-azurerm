package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

var kubernetesNodePoolDataSourceTests = map[string]func(t *testing.T){
	"basic": testAccAzureRMKubernetesClusterNodePoolDataSource_basic,
}

func TestAccAzureRMKubernetesClusterNodePoolDataSource_basic(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesClusterNodePoolDataSource_basic(t)
}

func testAccAzureRMKubernetesClusterNodePoolDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster_node_pool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesClusterNodePoolDataSource_basicConfig(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "node_count", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "Staging"),
				),
			},
		},
	})
}

func testAccAzureRMKubernetesClusterNodePoolDataSource_basicConfig(data acceptance.TestData) string {
	template := testAccAzureRMKubernetesClusterNodePool_manualScaleConfig(data)
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster_node_pool" "test" {
  name                    = azurerm_kubernetes_cluster_node_pool.test.name
  kubernetes_cluster_name = azurerm_kubernetes_cluster.test.name
  resource_group_name     = azurerm_kubernetes_cluster.test.resource_group_name
}
`, template)
}
