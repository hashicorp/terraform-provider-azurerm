// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package containers_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type KubernetesAutomaticClusterDataSource struct{}

func TestAccDataSourceKubernetesAutomaticCluster_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basicConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("kubelet_identity.0.object_id").Exists(),
				check.That(data.ResourceName).Key("kubelet_identity.0.client_id").Exists(),
				check.That(data.ResourceName).Key("kubelet_identity.0.user_assigned_identity_id").Exists(),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.principal_id").Exists(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").Exists(),
				check.That(data.ResourceName).Key("oidc_issuer_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("oidc_issuer_url").Exists(),
				check.That(data.ResourceName).Key("role_based_access_control_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("azure_active_directory_role_based_access_control.0.azure_rbac_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("azure_active_directory_role_based_access_control.0.tenant_id").Exists(),
				check.That(data.ResourceName).Key("azure_policy_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("network.0.network_plugin").HasValue("azure"),
				check.That(data.ResourceName).Key("network.0.outbound_type").Exists(),
				check.That(data.ResourceName).Key("storage.0.disk_driver_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("storage.0.file_driver_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("storage.0.snapshot_controller_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("network.0.network_policy").Exists(),
				check.That(data.ResourceName).Key("network.0.load_balancer_sku").Exists(),
				check.That(data.ResourceName).Key("network.0.service_cidr").Exists(),
				check.That(data.ResourceName).Key("network.0.dns_service_ip").Exists(),
				check.That(data.ResourceName).Key("agent_pools.#").Exists(),
				check.That(data.ResourceName).Key("agent_pools.0.name").Exists(),
				check.That(data.ResourceName).Key("agent_pools.0.type").Exists(),
				check.That(data.ResourceName).Key("agent_pools.0.vm_size").Exists(),
				check.That(data.ResourceName).Key("agent_pools.0.os_type").Exists(),
				check.That(data.ResourceName).Key("agent_pools.0.orchestrator_version").Exists(),
			),
		},
	})
}

func TestAccDataSourceKubernetesAutomaticCluster_serviceMesh(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.serviceMesh(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("service_mesh.#").HasValue("1"),
				check.That(data.ResourceName).Key("service_mesh.0.internal_ingress_gateway_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("service_mesh.0.external_ingress_gateway_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("service_mesh.0.revisions.0").HasValue("asm-1-29"),
			),
		},
	})
}

func TestAccDataSourceKubernetesAutomaticCluster_apiServerAuthorizedIPRanges(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.apiServerAuthorizedIPRanges(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("api_server_access.#").HasValue("1"),
				check.That(data.ResourceName).Key("api_server_access.0.authorized_ip_ranges.#").HasValue("2"),
			),
		},
	})
}

func TestAccDataSourceKubernetesAutomaticCluster_privateClusterEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.privateClusterEnabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("private_cluster.#").HasValue("1"),
				check.That(data.ResourceName).Key("private_cluster.0.public_fully_qualified_domain_name_enabled").HasValue("true"),
			),
		},
	})
}

func (KubernetesAutomaticClusterDataSource) basicConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_automatic_cluster" "test" {
  name                = azurerm_kubernetes_automatic_cluster.test.name
  resource_group_name = azurerm_kubernetes_automatic_cluster.test.resource_group_name
}
`, KubernetesAutomaticClusterResource{}.basic(data))
}

func (KubernetesAutomaticClusterDataSource) serviceMesh(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%[1]d"
  location = "%[2]s"
}

%[3]s

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  hosted_system {
    node_subnet_id        = azurerm_subnet.node.id
    system_node_subnet_id = azurerm_subnet.systemnode.id
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  api_server_access {
    subnet_id = azurerm_subnet.api.id
  }

  service_mesh {
    internal_ingress_gateway_enabled = true
    external_ingress_gateway_enabled = true
    revisions                        = ["asm-1-29"]
  }
}

data "azurerm_kubernetes_automatic_cluster" "test" {
  name                = azurerm_kubernetes_automatic_cluster.test.name
  resource_group_name = azurerm_kubernetes_automatic_cluster.test.resource_group_name
}
`, data.RandomInteger, data.Locations.Primary, KubernetesAutomaticClusterResource{}.networkTemplate(data))
}

func (KubernetesAutomaticClusterDataSource) apiServerAuthorizedIPRanges(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_automatic_cluster" "test" {
  name                = azurerm_kubernetes_automatic_cluster.test.name
  resource_group_name = azurerm_kubernetes_automatic_cluster.test.resource_group_name
}
`, KubernetesAutomaticClusterResource{}.apiServerAuthorizedIPRangesConfig(data))
}

func (KubernetesAutomaticClusterDataSource) privateClusterEnabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_automatic_cluster" "test" {
  name                = azurerm_kubernetes_automatic_cluster.test.name
  resource_group_name = azurerm_kubernetes_automatic_cluster.test.resource_group_name
}
`, KubernetesAutomaticClusterResource{}.privateClusterWithPublicFQDNConfig(data))
}
