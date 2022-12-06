package containers_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2022-09-02-preview/agentpools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2022-09-02-preview/managedclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type KubernetesClusterResource struct{}

var (
	olderKubernetesVersion        = "1.23.12"
	currentKubernetesVersion      = "1.24.3"
	olderKubernetesVersionAlias   = "1.23"
	currentKubernetesVersionAlias = "1.24"
)

func TestAccKubernetesCluster_hostEncryption(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.hostEncryption(data, currentKubernetesVersion),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("default_node_pool.0.enable_host_encryption").HasValue("true"),
			),
		},
	})
}

func TestAccKubernetesCluster_dedicatedHost(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dedicatedHost(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccKubernetesCluster_runCommand(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.runCommand(data, currentKubernetesVersion, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("run_command_enabled").HasValue("true"),
			),
		},
		{
			Config: r.runCommand(data, currentKubernetesVersion, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("run_command_enabled").HasValue("false"),
			),
		},
	})
}

func TestAccKubernetesCluster_workloadAutoscalerProfileKedaOnOff(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.workloadAutoscalerProfileKeda(data, currentKubernetesVersion, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("workload_autoscaler_profile.0.keda_enabled").HasValue("true"),
			),
		},
		{
			Config: r.workloadAutoscalerProfileKeda(data, currentKubernetesVersion, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("workload_autoscaler_profile.0.keda_enabled").HasValue("false"),
			),
		},
	})
}

func TestAccKubernetesCluster_workloadAutoscalerProfileKedaOnAbsent(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.workloadAutoscalerProfileKeda(data, currentKubernetesVersion, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("workload_autoscaler_profile.0.keda_enabled").HasValue("true"),
			),
		},
		{
			Config: r.basicVMSSConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesCluster_edgeZone(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.edgeZone(data, currentKubernetesVersion, "Test1"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.edgeZone(data, currentKubernetesVersion, "Test2"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func (t KubernetesClusterResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := managedclusters.ParseManagedClusterID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Containers.KubernetesClustersClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading Kubernetes Cluster (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.Model != nil && resp.Model.Id != nil), nil
}

func (KubernetesClusterResource) updateDefaultNodePoolAgentCount(nodeCount int) acceptance.ClientCheckFunc {
	return func(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) error {
		nodePoolName := state.Attributes["default_node_pool.0.name"]
		clusterName := state.Attributes["name"]
		resourceGroup := state.Attributes["resource_group_name"]

		agentPoolId := agentpools.NewAgentPoolID(clients.Account.SubscriptionId, resourceGroup, clusterName, nodePoolName)
		nodePool, err := clients.Containers.AgentPoolsClient.Get(ctx, agentPoolId)
		if err != nil {
			return fmt.Errorf("Bad: Get on agentPoolsClient: %+v", err)
		}

		if nodePool.HttpResponse.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Node Pool %q (Kubernetes Cluster %q / Resource Group: %q) does not exist", nodePoolName, clusterName, resourceGroup)
		}

		if nodePool.Model == nil || nodePool.Model.Properties == nil {
			return fmt.Errorf("Bad: Node Pool %q (Kubernetes Cluster %q / Resource Group: %q): `properties` was nil", nodePoolName, clusterName, resourceGroup)
		}

		nodePool.Model.Properties.Count = utils.Int64(int64(nodeCount))

		future, err := clients.Containers.AgentPoolsClient.CreateOrUpdate(ctx, agentPoolId, *nodePool.Model)
		if err != nil {
			return fmt.Errorf("Bad: updating node pool %q: %+v", nodePoolName, err)
		}

		if err := future.Poller.PollUntilDone(); err != nil {
			return fmt.Errorf("Bad: waiting for update of node pool %q: %+v", nodePoolName, err)
		}

		return nil
	}
}

func (KubernetesClusterResource) hostEncryption(data acceptance.TestData, controlPlaneVersion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"
  kubernetes_version  = %q

  default_node_pool {
    name                   = "default"
    node_count             = 1
    vm_size                = "Standard_DS2_v2"
    enable_host_encryption = true
  }

  identity {
    type = "SystemAssigned"
  }
}
  `, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, controlPlaneVersion)
}

func (KubernetesClusterResource) dedicatedHost(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%[1]d"
  location = "%[2]s"
}

resource "azurerm_dedicated_host_group" "test" {
  name                        = "acctestDHG-compute-%[1]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  platform_fault_domain_count = 3
  automatic_placement_enabled = true
}

resource "azurerm_dedicated_host" "test" {
  name                    = "acctest-DH-%[1]d"
  location                = azurerm_resource_group.test.location
  dedicated_host_group_id = azurerm_dedicated_host_group.test.id
  sku_name                = "DSv3-Type1"
  platform_fault_domain   = 0
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_resource_group.test.id
  principal_id         = azurerm_user_assigned_identity.test.principal_id
  role_definition_name = "Contributor"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[1]d"

  default_node_pool {
    name          = "default"
    node_count    = 1
    vm_size       = "Standard_D2s_v3"
    host_group_id = azurerm_dedicated_host_group.test.id
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  depends_on = [
    azurerm_role_assignment.test,
    azurerm_dedicated_host.test
  ]
}
  `, data.RandomInteger, data.Locations.Primary)
}

func (KubernetesClusterResource) runCommand(data acceptance.TestData, controlPlaneVersion string, runCommandEnabled bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"
  kubernetes_version  = %q
  run_command_enabled = %t

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_DS2_v2"
  }

  identity {
    type = "SystemAssigned"
  }
}
  `, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, controlPlaneVersion, runCommandEnabled)
}

func (KubernetesClusterResource) workloadAutoscalerProfileKeda(data acceptance.TestData, controlPlaneVersion string, kedaEnabled bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"
  kubernetes_version  = %q

  workload_autoscaler_profile {
    keda_enabled = %t
  }

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_DS2_v2"
  }

  identity {
    type = "SystemAssigned"
  }
}
  `, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, controlPlaneVersion, kedaEnabled)
}

func (r KubernetesClusterResource) upgradeSettingsConfig(data acceptance.TestData, maxSurge string) string {
	if maxSurge != "" {
		maxSurge = fmt.Sprintf(`upgrade_settings {
    max_surge = %q
  }`, maxSurge)
	}

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_DS2_v2"
    %s
  }

  identity {
    type = "SystemAssigned"
  }
}
  `, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, maxSurge)
}

func TestAccResourceKubernetesCluster_roleBasedAccessControlAAD_VOneDotTwoFourDotThree(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.roleBasedAccessControlAADManagedConfigVOneDotTwoFourDotThree(data, ""),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("kube_config.#").HasValue("1"),
				check.That(data.ResourceName).Key("kube_config.0.host").IsSet(),
			),
		},
	})
}

func (KubernetesClusterResource) edgeZone(data acceptance.TestData, controlPlaneVersion, tag string) string {
	// WestUS has an edge zone available - so hard-code to that
	data.Locations.Primary = "westus"

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}
data "azurerm_extended_locations" "test" {
  location = azurerm_resource_group.test.location
}
resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"
  kubernetes_version  = %q
  run_command_enabled = true
  edge_zone           = data.azurerm_extended_locations.test.extended_locations[0]
  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_DS2_v2"
  }
  identity {
    type = "SystemAssigned"
  }
  tags = {
    ENV = "%s"
  }
}
  `, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, controlPlaneVersion, tag)
}
