package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

var (
	olderKubernetesVersion   = "1.18.14"
	currentKubernetesVersion = "1.19.6"
)

func TestAccAzureRMKubernetes_all(t *testing.T) {
	// NOTE: this test is no longer used, but this assignment kicks around temporarily
	// to allow us to migrate off this without causing conflicts in open PR's
	_ = map[string]map[string]func(t *testing.T){
		"auth":               kubernetesAuthTests,
		"clusterAddOn":       kubernetesAddOnTests,
		"datasource":         kubernetesDataSourceTests,
		"network":            kubernetesNetworkAuthTests,
		"nodePool":           kubernetesNodePoolTests,
		"nodePoolDataSource": kubernetesNodePoolDataSourceTests,
		"other":              kubernetesOtherTests,
		"scaling":            kubernetesScalingTests,
		"upgrade":            kubernetesUpgradeTests,
	}

	t.Skip("Skipping since this is being run Individually")
}

func testCheckAzureRMKubernetesClusterExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Containers.KubernetesClustersClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Managed Kubernetes Cluster: %s", name)
		}

		aks, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on kubernetesClustersClient: %+v", err)
		}

		if aks.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Managed Kubernetes Cluster %q (Resource Group: %q) does not exist", name, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMKubernetesClusterDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).Containers.KubernetesClustersClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_kubernetes_cluster" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, name)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Managed Kubernetes Cluster still exists:\n%#v", resp)
		}
	}

	return nil
}

func kubernetesClusterUpdateNodePoolCount(resourceName string, nodeCount int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Containers.AgentPoolsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		nodePoolName := rs.Primary.Attributes["default_node_pool.0.name"]
		clusterName := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		nodePool, err := client.Get(ctx, resourceGroup, clusterName, nodePoolName)
		if err != nil {
			return fmt.Errorf("Bad: Get on agentPoolsClient: %+v", err)
		}

		if nodePool.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Node Pool %q (Kubernetes Cluster %q / Resource Group: %q) does not exist", nodePoolName, clusterName, resourceGroup)
		}

		if nodePool.ManagedClusterAgentPoolProfileProperties == nil {
			return fmt.Errorf("Bad: Node Pool %q (Kubernetes Cluster %q / Resource Group: %q): `properties` was nil", nodePoolName, clusterName, resourceGroup)
		}

		nodePool.ManagedClusterAgentPoolProfileProperties.Count = utils.Int32(int32(nodeCount))

		future, err := client.CreateOrUpdate(ctx, resourceGroup, clusterName, nodePoolName, nodePool)
		if err != nil {
			return fmt.Errorf("Bad: updating node pool %q: %+v", nodePoolName, err)
		}

		if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Bad: waiting for update of node pool %q: %+v", nodePoolName, err)
		}

		return nil
	}
}
