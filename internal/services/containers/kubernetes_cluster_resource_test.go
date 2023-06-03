package containers_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-04-02-preview/agentpools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type KubernetesClusterResource struct{}

var (
	olderKubernetesVersion        = "1.24.9"
	currentKubernetesVersion      = "1.25.5"
	olderKubernetesVersionAlias   = "1.24"
	currentKubernetesVersionAlias = "1.25"
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

func TestAccKubernetesCluster_keyVaultKms(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.azureKeyVaultKms(data, currentKubernetesVersion, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.azureKeyVaultKms(data, currentKubernetesVersion, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccKubernetesCluster_storageProfile(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.storageProfile(data, currentKubernetesVersion),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccKubernetesCluster_workloadAutoscalerProfileKedaToggle(t *testing.T) {
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
		data.ImportStep(),
		{
			Config: r.workloadAutoscalerProfileKeda(data, currentKubernetesVersion, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("workload_autoscaler_profile.0.keda_enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesCluster_imageCleanerSecurityProfileToggle(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.imageCleanerSecurityProfile(data, currentKubernetesVersion, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("image_cleaner_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("image_cleaner_interval_hours").HasValue("96"),
			),
		},
		{
			Config: r.imageCleanerSecurityProfile(data, currentKubernetesVersion, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("image_cleaner_enabled").HasValue("false"),
			),
		},
	})
}

func TestAccKubernetesCluster_workloadAutoscalerProfileVerticalPodAutoscalerToggle(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.workloadAutoscalerProfileVerticalPodAutoscaler(data, currentKubernetesVersion, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.workloadAutoscalerProfileVerticalPodAutoscaler(data, currentKubernetesVersion, false),
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
	id, err := commonids.ParseKubernetesClusterID(state.ID)
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

		if response.WasNotFound(nodePool.HttpResponse) {
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

func TestAccKubernetesCluster_dnsPrefix(t *testing.T) {
	// regression test case for issue #20806
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	dnsPrefix := fmt.Sprintf("1stCluster%d", data.RandomInteger)

	r := KubernetesClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dnsPrefix(data, currentKubernetesVersion),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("dns_prefix").HasValue(dnsPrefix),
			),
		},
	})
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
  sku_name                = "FSv2-Type2"
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

func (KubernetesClusterResource) workloadAutoscalerProfileVerticalPodAutoscaler(data acceptance.TestData, controlPlaneVersion string, enabled bool) string {
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
    vertical_pod_autoscaler_enabled = %t
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
  `, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, controlPlaneVersion, enabled)
}

func (KubernetesClusterResource) imageCleanerSecurityProfile(data acceptance.TestData, controlPlaneVersion string, enabled bool) string {
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

  image_cleaner_enabled        = %t
  image_cleaner_interval_hours = 96

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_DS2_v2"
  }

  identity {
    type = "SystemAssigned"
  }
}
  `, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, controlPlaneVersion, enabled)
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

func TestAccResourceKubernetesCluster_roleBasedAccessControlAAD_VOneDotTwoFourDotNine(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.roleBasedAccessControlAADManagedConfigVOneDotTwoFourDotNine(data, ""),
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

func (KubernetesClusterResource) azureKeyVaultKms(data acceptance.TestData, controlPlaneVersion string, enabled bool) string {
	kmsBlock := ""
	if enabled {
		kmsBlock = `
  key_management_service {
    key_vault_key_id = azurerm_key_vault_key.test.id
  }`
	}

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%[1]d"
  location = "%[2]s"
}

resource "azurerm_key_vault" "test" {
  name                      = substr("acctest%[1]d", 0, 24)
  location                  = azurerm_resource_group.test.location
  resource_group_name       = azurerm_resource_group.test.name
  tenant_id                 = data.azurerm_client_config.current.tenant_id
  enable_rbac_authorization = true
  sku_name                  = "standard"
}

resource "azurerm_role_assignment" "test_admin" {
  scope                = azurerm_key_vault.test.id
  role_definition_name = "Key Vault Administrator"
  principal_id         = data.azurerm_client_config.current.object_id
}

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_key_vault.test.id
  role_definition_name = "Key Vault Crypto User"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_key_vault_key" "test" {
  name         = "etcd-encryption"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]

  depends_on = [azurerm_role_assignment.test_admin]
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  node_resource_group = "${azurerm_resource_group.test.name}-infra"
  dns_prefix          = "acctestaks%[1]d"
  kubernetes_version  = %[3]q

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_DS2_v2"
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
  %[4]s
}
`, data.RandomInteger, data.Locations.Primary, controlPlaneVersion, kmsBlock)
}

func (KubernetesClusterResource) storageProfile(data acceptance.TestData, controlPlaneVersion string) string {
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
    name       = "default"
    node_count = 1
    vm_size    = "Standard_DS2_v2"
  }

  identity {
    type = "SystemAssigned"
  }

  storage_profile {
    blob_driver_enabled         = true
    disk_driver_enabled         = true
    disk_driver_version         = "v1"
    file_driver_enabled         = false
    snapshot_controller_enabled = false
  }
}
  `, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, controlPlaneVersion)
}

func (KubernetesClusterResource) dnsPrefix(data acceptance.TestData, controlPlaneVersion string) string {
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
  dns_prefix          = "1stCluster%d"
  kubernetes_version  = %q

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_DS2_v2"
  }

  identity {
    type = "SystemAssigned"
  }
}
  `, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, controlPlaneVersion)
}

func (KubernetesClusterResource) customCaTrustCertificates(data acceptance.TestData, fakeCertsList []string) string {
	return fmt.Sprintf(`
%s
data "azurerm_kubernetes_cluster" "test" {
  name                = azurerm_kubernetes_cluster.test.name
  resource_group_name = azurerm_kubernetes_cluster.test.resource_group_name
}
`, KubernetesClusterResource{}.customCATrustCertificates(data, fakeCertsList))
}

func TestAccKubernetesCluster_customCaTrustCerts(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}

	fakeCertList := []string{
		"LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURjRENDQWxpZ0F3SUJBZ0lFU1QwSUhEQU5CZ2txaGtpRzl3MEJBUXNGQURCUk1Rc3dDUVlEVlFRR0V3SlEKVERFTk1Bc0dBMVVFQXd3RVZHVnpkREVWTUJNR0ExVUVCd3dNUkdWbVlYVnNkQ0JEYVhSNU1Sd3dHZ1lEVlFRSwpEQk5FWldaaGRXeDBJRU52YlhCaGJua2dUSFJrTUI0WERUSXpNRFV5T0RFeE1qY3dNMW9YRFRNek1EVXlOVEV4Ck1qY3dNMW93VVRFTE1Ba0dBMVVFQmhNQ1VFd3hEVEFMQmdOVkJBTU1CRlJsYzNReEZUQVRCZ05WQkFjTURFUmwKWm1GMWJIUWdRMmwwZVRFY01Cb0dBMVVFQ2d3VFJHVm1ZWFZzZENCRGIyMXdZVzU1SUV4MFpEQ0NBU0l3RFFZSgpLb1pJaHZjTkFRRUJCUUFEZ2dFUEFEQ0NBUW9DZ2dFQkFLN2JIYWtxSkdRMWVBOUFHUmlhNGl2anNDRXlGMDhDCjNpSzJZeWthNkREeldmTk1tRWpOUjJiQVZOMEhlLy9pWTd1VjJ2dXl6V1UxMzZGVkdMZkdyeTZGOHNQQUZaSzYKSE4vcWk1QVp6MUpoOGdWSTRwS1pjZEFxQS81clF3VVlvWVN3Q245dGVOYytsbU1ZUk5OcTVwdlV2NjcrNEM3MgpPc3BOSUxSclhBbWNUb1YveVRZVzFKWDBOeEJJSHZZaFZXUE9LQXpRZDQ5UEpSeFpqMUgydCszMEFsazgzTDFwClFzTGx2SzV3MjJpeXdkYVpRN1lmV0xXd1hPQzVPWXdRTUw1R3BHUFNQaEdxdjhqSUhpcHBVeTdrRDlNWFFZOFoKdDl2QkczMzVWSEdlUjI2QnNQQXRFbTJjR05ocjA5cmRvdWJGd2tDR05OYXNVamFoVW9CKzhPY0NBd0VBQWFOUQpNRTR3SFFZRFZSME9CQllFRk9CNmNpTGtUL21Cc2xXSm5Na2phQzZqbjd4ek1COEdBMVVkSXdRWU1CYUFGT0I2CmNpTGtUL21Cc2xXSm5Na2phQzZqbjd4ek1Bd0dBMVVkRXdRRk1BTUJBZjh3RFFZSktvWklodmNOQVFFTEJRQUQKZ2dFQkFKTklHdHJpeFlCRUc1Yy9iQWdOMHlMOEJvOW9nN29ha0hVMUc5TjBxOUNWWXhjOVhma2ZUaEhYOVBUeApMbVNGcHJEQlAyYnVGTzVIUDFpbnNFT1E2N1lGanAvRjVJWGdaQ2twZUpGdDBTL0R3N2ZRbFJJN2RCNGQzNmIzCmE1R2txU0M4aFlZemxLUm9DRGNhalp4QmdoVUFxK0tnTnV4RmNsM1Fnd1Uyam1QbkU4a1A4TmgyM3hlVUJ3WEkKL3pqbU1rdjV4SFhKdHBpdlpzTlpSSUttQW56RU9TWGlRK2JMTStTdlhtSkhYd29YYTZyTXg4YmkySzV4WkhIRwpkUHA1TnQ3L2dxOUdXcm95SkVjSFpEclBiSnR2WGFibTZYUXpxTTFYUzA3SDlaSFBXc0dENGlBM1k0T3JUUlRCClZ5blRPUDl5U3cwbklaVEk4YjZuR2RHTzBOOD0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQ==",
		"LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURmakNDQW1hZ0F3SUJBZ0lFZnlWdk56QU5CZ2txaGtpRzl3MEJBUXNGQURCWU1Rc3dDUVlEVlFRR0V3SlEKVERFVU1CSUdBMVVFQXd3TFJtRnJaU0JEWlhKMElESXhGVEFUQmdOVkJBY01ERVJsWm1GMWJIUWdRMmwwZVRFYwpNQm9HQTFVRUNnd1RSR1ZtWVhWc2RDQkRiMjF3WVc1NUlFeDBaREFlRncweU16QTJNRFF3TnpJME1qZGFGdzB5Ck5UQTJNRE13TnpJME1qZGFNRmd4Q3pBSkJnTlZCQVlUQWxCTU1SUXdFZ1lEVlFRRERBdEdZV3RsSUVObGNuUWcKTWpFVk1CTUdBMVVFQnd3TVJHVm1ZWFZzZENCRGFYUjVNUnd3R2dZRFZRUUtEQk5FWldaaGRXeDBJRU52YlhCaApibmtnVEhSa01JSUJJakFOQmdrcWhraUc5dzBCQVFFRkFBT0NBUThBTUlJQkNnS0NBUUVBMENTdVdUaGNjSG5MCkhFdjk4SUVNc2JLY3h4YVh4YTZiRXl1Yy9sUjRackpVN2p6eVlWNGVscTV5WTgwdDFCM0MyV3E2SXFoajErSGYKYW0xaStsU1FTejM1eWNnTWlwSWp2cUxKOVIzMVF0Wi9TRURkdGV2b2JqbytEa1dCOE55cG9Ia0pVbEIyQnR6ZgpOK09KeVFSdXU1b1cya2c5OE5Bd3JuTGpmQ0lremVWcFh5d0l4Tkx2ZmFrVGxpNWpYdG9WWG5pOTU5bmtINWVwClkrRnVoSEQwaU5CS25XYVkxR2QwVGhhSHNwTERmNFUycmo2WE5SZHd6QVZoVkdhUm02cndvSHRZeDVrYys1ZWMKQ0F4UEdRWFRzTzJUTHVrQzJ2YXI0M3RUM0ZjSC9taDRST2JaaThZS2xSQ3Fldm1QU1RmZ293RUFkTjlvSmxyRApXN2lzN2NnQjhRSURBUUFCbzFBd1RqQWRCZ05WSFE0RUZnUVVuRkRqN0pBQW9WZ2NzQkgyNzdMOHZlM0Q4U293Ckh3WURWUjBqQkJnd0ZvQVVuRkRqN0pBQW9WZ2NzQkgyNzdMOHZlM0Q4U293REFZRFZSMFRCQVV3QXdFQi96QU4KQmdrcWhraUc5dzBCQVFzRkFBT0NBUUVBT0diT0Zyek4rN2YxbzhJSDNtMXZxT3IyTUtvNEZMWExGRjBVbEhkNApwZXRhL05aQjArUmQ3TnUrOCtnUnlUbEJWZU9EZjN5SXU0TlFCUU92MlNqdS9Jakd0MUtmaUF3WkUwT1RUQXc3CnhIWStsMVBJWEFFVWNqNk00cjFKQzc4ZVZrc2pycTZoV1RPZ0RrSVZuRjY3bXlReXduR25EY1k0d0Fqc2pUajgKKzR4NTIrRi9QaVNQVGtjUFNuN0s2UjQzaEt5QUs2Z0poOHE5cVNhME5RQ2U2czhwTGU2SVY5SElWVVFFVERVOQpsM1VWWHNBMGx4dlB0blU1TXo2QWQ5cDA5L2w4d3o0cUdBdGFCUEd3K0R2cTNlaHdTd2VZZ3VHSktDQjhjb01JCjJRVUo0Zi9mNkFNVWtMeWxYZ3RSUEt1QjA3d3YwTmk1eWI5MjlFY1FJQ0l2dFE9PQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0t",
	}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.customCaTrustCertificates(data, fakeCertList),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("custom_ca_trust_certificates.0").Exists(),
				check.That(data.ResourceName).Key("custom_ca_trust_certificates.1").Exists(),
			),
		},
		{
			Config: r.customCaTrustCertificates(data, nil),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("custom_ca_trust_certificates.0").DoesNotExist(),
				check.That(data.ResourceName).Key("custom_ca_trust_certificates.1").DoesNotExist(),
			),
		},
	})
}
