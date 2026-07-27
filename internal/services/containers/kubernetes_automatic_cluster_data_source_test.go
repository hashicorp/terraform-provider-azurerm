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
				check.That(data.ResourceName).Key("kube_admin_config.#").HasValue("0"),
				check.That(data.ResourceName).Key("kube_admin_config_raw").HasValue(""),
				check.That(data.ResourceName).Key("kubelet_identity.0.object_id").Exists(),
				check.That(data.ResourceName).Key("kubelet_identity.0.client_id").Exists(),
				check.That(data.ResourceName).Key("kubelet_identity.0.user_assigned_identity_id").Exists(),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.principal_id").Exists(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").Exists(),
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
				check.That(data.ResourceName).Key("service_mesh.0.mode").HasValue("Istio"),
				check.That(data.ResourceName).Key("service_mesh.0.internal_ingress_gateway_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("service_mesh.0.external_ingress_gateway_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("service_mesh.0.revisions.0").HasValue("asm-1-28"),
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
				check.That(data.ResourceName).Key("api_server_authorized_ip_ranges.#").HasValue("2"),
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
				check.That(data.ResourceName).Key("private_cluster_enabled").HasValue("true"),
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
%s

data "azurerm_kubernetes_automatic_cluster" "test" {
  name                = azurerm_kubernetes_automatic_cluster.test.name
  resource_group_name = azurerm_kubernetes_automatic_cluster.test.resource_group_name
}
`, KubernetesAutomaticClusterResource{}.serviceMeshProfile(data, true, true))
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
