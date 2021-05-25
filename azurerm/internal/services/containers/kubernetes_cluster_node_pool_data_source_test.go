package containers_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type KubernetesClusterNodePoolDataSource struct {
}

var kubernetesNodePoolDataSourceTests = map[string]func(t *testing.T){
	"basic": testAccKubernetesClusterNodePoolDataSource_basic,
}

func TestAccKubernetesClusterNodePoolDataSource_basic(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccKubernetesClusterNodePoolDataSource_basic(t)
}

func testAccKubernetesClusterNodePoolDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basicConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("node_count").HasValue("1"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("Staging"),
			),
		},
	})
}

func (KubernetesClusterNodePoolDataSource) basicConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster_node_pool" "test" {
  name                    = azurerm_kubernetes_cluster_node_pool.test.name
  kubernetes_cluster_name = azurerm_kubernetes_cluster.test.name
  resource_group_name     = azurerm_kubernetes_cluster.test.resource_group_name
}
`, KubernetesClusterNodePoolResource{}.manualScaleConfig(data))
}
