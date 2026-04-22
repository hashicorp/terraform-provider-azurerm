// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package containers_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-10-01/agentpools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type KubernetesAutomaticClusterResource struct{}

var (
	olderKubernetesAutomaticVersion        = "1.33.3"
	currentKubernetesAutomaticVersion      = "1.34.2"
	olderKubernetesAutomaticVersionAlias   = "1.33"
	currentKubernetesAutomaticVersionAlias = "1.34"
)

func TestAccKubernetesAutomaticCluster_automaticSKU(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.automaticSKU(data, currentKubernetesAutomaticVersion),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccKubernetesAutomaticCluster_hostEncryption(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.hostEncryption(data, currentKubernetesAutomaticVersion),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

// func TestAccKubernetesAutomaticCluster_dedicatedHost(t *testing.T) {
//	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
//	r := KubernetesAutomaticClusterResource{}
//
//	data.ResourceTest(t, r, []acceptance.TestStep{
//		{
//			Config: r.dedicatedHost(data),
//			Check: acceptance.ComposeTestCheckFunc(
//				check.That(data.ResourceName).ExistsInAzure(r),
//			),
//		},
//	})
//}

func TestAccKubernetesAutomaticCluster_runCommand(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.runCommand(data, currentKubernetesAutomaticVersion, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("run_command_enabled").HasValue("true"),
			),
		},
		{
			Config: r.runCommand(data, currentKubernetesAutomaticVersion, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("run_command_enabled").HasValue("false"),
			),
		},
	})
}

func TestAccKubernetesAutomaticCluster_keyVaultKms(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.azureKeyVaultKms(data, currentKubernetesAutomaticVersion, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.azureKeyVaultKms(data, currentKubernetesAutomaticVersion, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccKubernetesAutomaticCluster_storageProfile(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.storageProfile(data, currentKubernetesAutomaticVersion),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccKubernetesAutomaticCluster_nodeProvisioningProfileUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.nodeProvisioningProfile(data, "Auto"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.nodeProvisioningProfile(data, "None"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.nodeProvisioningProfile(data, "Auto"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.nodeProvisioningProfileRemoved(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

// func TestAccKubernetesAutomaticCluster_edgeZone(t *testing.T) {
//	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
//	r := KubernetesAutomaticClusterResource{}
//
//	data.ResourceTest(t, r, []acceptance.TestStep{
//		{
//			Config: r.edgeZone(data, currentKubernetesAutomaticVersion, "Test1"),
//			Check: acceptance.ComposeTestCheckFunc(
//				check.That(data.ResourceName).ExistsInAzure(r),
//			),
//		},
//		{
//			Config: r.edgeZone(data, currentKubernetesAutomaticVersion, "Test2"),
//			Check: acceptance.ComposeTestCheckFunc(
//				check.That(data.ResourceName).ExistsInAzure(r),
//			),
//		},
//	})
//}

func TestAccKubernetesAutomaticCluster_bootstrapProfile(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.networkIsolatedBootstrapProfileArtifactSourceDirect(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.networkIsolatedBootstrapProfileArtifactSourceCache(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.networkIsolatedBootstrapProfileRemoved(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_upgradeOverrideSetting(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.upgradeOverrideSetting(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.upgradeOverrideSetting(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.upgradeOverrideSetting(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t KubernetesAutomaticClusterResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commonids.ParseKubernetesClusterIDInsensitively(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Containers.KubernetesClustersClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading Kubernetes Cluster (%s): %+v", id.String(), err)
	}

	return pointer.To(resp.Model != nil && resp.Model.Id != nil), nil
}

func (KubernetesAutomaticClusterResource) updateDefaultNodePoolAgentCount(nodeCount int) acceptance.ClientCheckFunc {
	return func(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) error {
		if _, ok := ctx.Deadline(); !ok {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, 1*time.Hour)
			defer cancel()
		}

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

		nodePool.Model.Properties.Count = pointer.To(int64(nodeCount))

		err = clients.Containers.AgentPoolsClient.CreateOrUpdateThenPoll(ctx, agentPoolId, *nodePool.Model, agentpools.DefaultCreateOrUpdateOperationOptions())
		if err != nil {
			return fmt.Errorf("Bad: updating node pool %q: %+v", nodePoolName, err)
		}

		return nil
	}
}

func TestAccKubernetesAutomaticCluster_dnsPrefix(t *testing.T) {
	// regression test case for issue #20806
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	dnsPrefix := fmt.Sprintf("1stCluster%d", data.RandomInteger)

	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dnsPrefix(data, currentKubernetesAutomaticVersion),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("dns_prefix").HasValue(dnsPrefix),
			),
		},
	})
}

func (KubernetesAutomaticClusterResource) automaticSKU(data acceptance.TestData, controlPlaneVersion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"
  kubernetes_version  = %q

  default_node_pool {
    name       = "default"
    node_count = 1
    //os_disk_type = "Ephemeral"
    //vm_size    = "standard_d4lds_v5"
  }

  identity {
    type = "SystemAssigned"
  }
}
  `, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, controlPlaneVersion)
}

func (KubernetesAutomaticClusterResource) hostEncryption(data acceptance.TestData, controlPlaneVersion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"
  kubernetes_version  = %q

  default_node_pool {
    name                    = "default"
    node_count              = 1
    host_encryption_enabled = true
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type = "SystemAssigned"
  }
}
  `, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, controlPlaneVersion)
}

// func (KubernetesAutomaticClusterResource) dedicatedHost(data acceptance.TestData) string {
//	return fmt.Sprintf(`
// provider "azurerm" {
//  features {}
// }
//
// resource "azurerm_resource_group" "test" {
//  name     = "acctestRG-aks-%[1]d"
//  location = "%[2]s"
// }
//
// resource "azurerm_virtual_network" "test" {
//  name                = "acctestvirtnet%[1]d"
//  address_space       = ["10.0.0.0/8"]
//  location            = azurerm_resource_group.test.location
//  resource_group_name = azurerm_resource_group.test.name
// }
//
// resource "azurerm_subnet" "test" {
//  name                 = "acctestsubnet%[1]d"
//  resource_group_name  = azurerm_resource_group.test.name
//  virtual_network_name = azurerm_virtual_network.test.name
//  address_prefixes     = ["10.1.0.0/16"]
//
//  delegation {
//    name = "aks-delegation"
//
//    service_delegation {
//      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
//      name    = "Microsoft.ContainerService/managedClusters"
//    }
//  }
// }
//
// resource "azurerm_subnet" "test1" {
//  name                 = "acctestsubnet1%[1]d"
//  resource_group_name  = azurerm_resource_group.test.name
//  virtual_network_name = azurerm_virtual_network.test.name
//  address_prefixes     = ["10.2.0.0/16"]
// }
//
// resource "azurerm_dedicated_host_group" "test" {
//  name                        = "acctestDHG-compute-%[1]d"
//  resource_group_name         = azurerm_resource_group.test.name
//  location                    = azurerm_resource_group.test.location
//  platform_fault_domain_count = 3
//  automatic_placement_enabled = true
// }
//
// resource "azurerm_dedicated_host" "test" {
//  name                    = "acctest-DH-%[1]d"
//  location                = azurerm_resource_group.test.location
//  dedicated_host_group_id = azurerm_dedicated_host_group.test.id
//  sku_name                = "FSv2-Type2"
//  platform_fault_domain   = 0
// }
//
// resource "azurerm_user_assigned_identity" "test" {
//  name                = "acctest%[2]s"
//  resource_group_name = azurerm_resource_group.test.name
//  location            = azurerm_resource_group.test.location
// }
//
// resource "azurerm_role_assignment" "test" {
//  scope                = azurerm_resource_group.test.id
//  principal_id         = azurerm_user_assigned_identity.test.principal_id
//  role_definition_name = "Contributor"
// }
//
// resource "azurerm_role_assignment" "test1" {
//  scope                = azurerm_subnet.test.id
//  role_definition_name = "Network Contributor"
//  principal_id         = azurerm_user_assigned_identity.test.principal_id
// }
//
// resource "azurerm_kubernetes_automatic_cluster" "test" {
//  name                = "acctestaks%[1]d"
//  location            = azurerm_resource_group.test.location
//  resource_group_name = azurerm_resource_group.test.name
//  dns_prefix          = "acctestaks%[1]d"
//
//  default_node_pool {
//    name           = "default"
//    node_count     = 1
//    vnet_subnet_id = azurerm_subnet.test1.id
//    host_group_id  = azurerm_dedicated_host_group.test.id
//    upgrade_settings {
//      max_surge = "10%%"
//    }
//  }
//
//  identity {
//    type         = "UserAssigned"
//    identity_ids = [azurerm_user_assigned_identity.test.id]
//  }
//
//  api_server_access_profile {
//    subnet_id = azurerm_subnet.test.id
//  }
//
//  depends_on = [
//    azurerm_dedicated_host.test
//  ]
// }
//  `, data.RandomInteger, data.Locations.Primary)
// }

func (KubernetesAutomaticClusterResource) runCommand(data acceptance.TestData, controlPlaneVersion string, runCommandEnabled bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"
  kubernetes_version  = %q
  run_command_enabled = %t

  default_node_pool {
    name       = "default"
    node_count = 1
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type = "SystemAssigned"
  }
}
  `, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, controlPlaneVersion, runCommandEnabled)
}

func (KubernetesAutomaticClusterResource) nodeProvisioningProfile(data acceptance.TestData, defaultNodePools string) string {
	return fmt.Sprintf(`provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%[1]d"
  location = "%[2]s"
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[1]d"
  node_provisioning_profile {
    default_node_pools = "%[3]s"
  }
  default_node_pool {
    name       = "default"
    node_count = 1
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, defaultNodePools)
}

func (KubernetesAutomaticClusterResource) nodeProvisioningProfileRemoved(data acceptance.TestData) string {
	return fmt.Sprintf(`provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%[1]d"
  location = "%[2]s"
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[1]d"
  default_node_pool {
    name       = "default"
    node_count = 1
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func TestAccResourceKubernetesAutomaticCluster_roleBasedAccessControlAAD_OlderKubernetesVersion(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.roleBasedAccessControlAADManagedConfigOlderKubernetesVersion(data, ""),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("kube_config.#").HasValue("1"),
				check.That(data.ResourceName).Key("kube_config.0.host").IsSet(),
			),
		},
	})
}

// func (KubernetesAutomaticClusterResource) edgeZone(data acceptance.TestData, controlPlaneVersion, tag string) string {
//	// WestUS has an edge zone available - so hard-code to that
//	data.Locations.Primary = "westus"
//
//	return fmt.Sprintf(`
// provider "azurerm" {
//  features {}
// }
// resource "azurerm_resource_group" "test" {
//  name     = "acctestRG-aks-%d"
//  location = "%s"
// }
// data "azurerm_extended_locations" "test" {
//  location = azurerm_resource_group.test.location
// }
// resource "azurerm_kubernetes_automatic_cluster" "test" {
//  name                = "acctestaks%d"
//  location            = azurerm_resource_group.test.location
//  resource_group_name = azurerm_resource_group.test.name
//  dns_prefix          = "acctestaks%d"
//  kubernetes_version  = %q
//  run_command_enabled = true
//  edge_zone           = data.azurerm_extended_locations.test.extended_locations[0]
//  default_node_pool {
//    name       = "default"
//    node_count = 1
//    vm_size    = "Standard_DS4_v2"
//    upgrade_settings {
//      max_surge = "10%%"
//    }
//  }
//  identity {
//    type = "SystemAssigned"
//  }
//  tags = {
//    ENV = "%s"
//  }
// }
// `, data.RandomInteger, data.Locations.Primary, data.RandomInteger, controlPlaneVersion, tag)
// }

func (KubernetesAutomaticClusterResource) azureKeyVaultKms(data acceptance.TestData, controlPlaneVersion string, enabled bool) string {
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


resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%[1]d"
  address_space       = ["10.0.0.0/8"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.1.0.0/16"]

  delegation {
    name = "aks-delegation"

    service_delegation {
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
      name    = "Microsoft.ContainerService/managedClusters"
    }
  }
}

resource "azurerm_subnet" "test1" {
  name                 = "acctestsubnet1%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.2.0.0/16"]
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

resource "azurerm_role_assignment" "test1" {
  scope                = azurerm_subnet.test.id
  role_definition_name = "Network Contributor"
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

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  node_resource_group = "${azurerm_resource_group.test.name}-infra"
  dns_prefix          = "acctestaks%[1]d"
  kubernetes_version  = %[3]q

  default_node_pool {
    name           = "default"
    node_count     = 1
    vnet_subnet_id = azurerm_subnet.test1.id
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  api_server_access_profile {
    subnet_id = azurerm_subnet.test.id
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
  %[4]s
}
`, data.RandomInteger, data.Locations.Primary, controlPlaneVersion, kmsBlock)
}

func (KubernetesAutomaticClusterResource) storageProfile(data acceptance.TestData, controlPlaneVersion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"
  kubernetes_version  = %q

  default_node_pool {
    name       = "default"
    node_count = 1
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type = "SystemAssigned"
  }

  storage_profile {
    blob_driver_enabled         = true
    disk_driver_enabled         = true
    file_driver_enabled         = false
    snapshot_controller_enabled = false
  }
}
  `, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, controlPlaneVersion)
}

func (KubernetesAutomaticClusterResource) dnsPrefix(data acceptance.TestData, controlPlaneVersion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "1stCluster%d"
  kubernetes_version  = %q

  default_node_pool {
    name       = "default"
    node_count = 1
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type = "SystemAssigned"
  }
}
  `, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, controlPlaneVersion)
}

func (KubernetesAutomaticClusterResource) upgradeOverrideSetting(data acceptance.TestData, isUpgradeOverrideSettingEnabled bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%[1]s"
  location = "%[2]s"
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%[1]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[1]s"

  default_node_pool {
    name       = "default"
    node_count = 1
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type = "SystemAssigned"
  }

  upgrade_override {
    effective_until       = "%[3]s"
    force_upgrade_enabled = %[4]t
  }
}
  `, data.RandomString, data.Locations.Primary, time.Now().UTC().Add(8*time.Minute).Format(time.RFC3339), isUpgradeOverrideSettingEnabled)
}

func (KubernetesAutomaticClusterResource) networkIsolatedBootstrapProfileTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%[2]d"
  location = "%[1]s"
}

resource "azurerm_container_registry" "registry" {
  name                          = "acctestacr%[2]d"
  resource_group_name           = azurerm_resource_group.test.name
  location                      = azurerm_resource_group.test.location
  sku                           = "Premium"
  public_network_access_enabled = false
  admin_enabled                 = false
  network_rule_bypass_option    = "None"
}

resource "azurerm_container_registry_cache_rule" "cache_rule" {
  name                  = "aks-managed-mcr"
  container_registry_id = azurerm_container_registry.registry.id
  target_repo           = "aks-managed-repository/*"
  source_repo           = "mcr.microsoft.com/*"
}

resource "azurerm_role_assignment" "aks_pull_from_acr" {
  scope                = azurerm_container_registry.registry.id
  role_definition_name = "AcrPull"
  principal_id         = azurerm_user_assigned_identity.aks_kubelet.principal_id
}

resource "azurerm_private_dns_zone" "acr_private_dns_zone" {
  name                = "privatelink.azurecr.io"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_private_dns_zone_virtual_network_link" "dns_vnet_link_acr" {
  name                  = "acctest-dns-vnet-link-acr"
  private_dns_zone_name = azurerm_private_dns_zone.acr_private_dns_zone.name
  resource_group_name   = azurerm_resource_group.test.name
  virtual_network_id    = azurerm_virtual_network.vnet.id
}

resource "azurerm_private_endpoint" "acr_private_endpoint" {
  name                = "acctest-acr-pe"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  subnet_id           = azurerm_subnet.vnet-nodepool.id

  private_service_connection {
    name                           = "acctest-sc-acr"
    private_connection_resource_id = azurerm_container_registry.registry.id
    is_manual_connection           = false
    subresource_names = [
      "registry"
    ]
  }

  private_dns_zone_group {
    name = "acctest-dns-group-acr"
    private_dns_zone_ids = [
      azurerm_private_dns_zone.acr_private_dns_zone.id
    ]
  }
}

resource "azurerm_virtual_network" "vnet" {
  name                = "acctestvnet"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["172.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%[2]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.vnet.name
  address_prefixes     = ["172.0.0.0/24"]

  delegation {
    name = "aks-delegation"

    service_delegation {
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
      name    = "Microsoft.ContainerService/managedClusters"
    }
  }
}

resource "azurerm_role_assignment" "test1" {
  scope                = azurerm_subnet.test.id
  role_definition_name = "Network Contributor"
  principal_id         = azurerm_user_assigned_identity.aks.principal_id
}

resource "azurerm_subnet" "vnet-nodepool" {
  name                 = "aks"
  virtual_network_name = azurerm_virtual_network.vnet.name
  resource_group_name  = azurerm_resource_group.test.name
  address_prefixes     = ["172.0.32.0/24"]
}

resource "azurerm_user_assigned_identity" "aks_kubelet" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "aks_kubelet-identity"
}

resource "azurerm_user_assigned_identity" "aks" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "aks-identity"
}

resource "azurerm_role_assignment" "aks_to_vnet" {
  scope                = azurerm_virtual_network.vnet.id
  role_definition_name = "Contributor"
  principal_id         = azurerm_user_assigned_identity.aks.principal_id
}

resource "azurerm_role_assignment" "aks_to_kubeletidentity" {
  scope                = azurerm_user_assigned_identity.aks_kubelet.id
  role_definition_name = "Managed Identity Operator"
  principal_id         = azurerm_user_assigned_identity.aks.principal_id
}


  `, data.Locations.Primary, data.RandomInteger)
}

func (r KubernetesAutomaticClusterResource) networkIsolatedBootstrapProfileArtifactSourceCache(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[2]d"

  default_node_pool {
    name       = "default"
    node_count = 1
    upgrade_settings {
      max_surge = "10%%"
    }
    vnet_subnet_id = azurerm_subnet.vnet-nodepool.id
  }

  network_profile {
    network_plugin_mode = "overlay"
    load_balancer_sku   = "standard"
    outbound_type       = "loadBalancer"
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.aks.id]
  }

  kubelet_identity {
    user_assigned_identity_id = azurerm_user_assigned_identity.aks_kubelet.id
    client_id                 = azurerm_user_assigned_identity.aks_kubelet.client_id
    object_id                 = azurerm_user_assigned_identity.aks_kubelet.principal_id
  }

  api_server_access_profile {
    subnet_id = azurerm_subnet.test.id
  }


  bootstrap_profile {
    artifact_source       = "Cache"
    container_registry_id = azurerm_container_registry.registry.id
  }
}`, r.networkIsolatedBootstrapProfileTemplate(data), data.RandomInteger)
}

func (r KubernetesAutomaticClusterResource) networkIsolatedBootstrapProfileArtifactSourceDirect(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[2]d"

  default_node_pool {
    name       = "default"
    node_count = 1
    upgrade_settings {
      max_surge = "10%%"
    }
    vnet_subnet_id = azurerm_subnet.vnet-nodepool.id
  }

  network_profile {
    network_plugin_mode = "overlay"
    load_balancer_sku   = "standard"
    outbound_type       = "loadBalancer"
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.aks.id]
  }

  kubelet_identity {
    user_assigned_identity_id = azurerm_user_assigned_identity.aks_kubelet.id
    client_id                 = azurerm_user_assigned_identity.aks_kubelet.client_id
    object_id                 = azurerm_user_assigned_identity.aks_kubelet.principal_id
  }

  api_server_access_profile {
    subnet_id = azurerm_subnet.test.id
  }

  bootstrap_profile {
    artifact_source = "Direct"
  }
}`, r.networkIsolatedBootstrapProfileTemplate(data), data.RandomInteger)
}

func (r KubernetesAutomaticClusterResource) networkIsolatedBootstrapProfileRemoved(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[2]d"

  default_node_pool {
    name       = "default"
    node_count = 1
    upgrade_settings {
      max_surge = "10%%"
    }
    vnet_subnet_id = azurerm_subnet.vnet-nodepool.id
  }

  network_profile {
    network_plugin_mode = "overlay"
    load_balancer_sku   = "standard"
    outbound_type       = "loadBalancer"
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.aks.id]
  }


  api_server_access_profile {
    subnet_id = azurerm_subnet.test.id
  }


  kubelet_identity {
    user_assigned_identity_id = azurerm_user_assigned_identity.aks_kubelet.id
    client_id                 = azurerm_user_assigned_identity.aks_kubelet.client_id
    object_id                 = azurerm_user_assigned_identity.aks_kubelet.principal_id
  }
}`, r.networkIsolatedBootstrapProfileTemplate(data), data.RandomInteger)
}
