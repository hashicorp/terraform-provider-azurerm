// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containers_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type KubernetesClusterDeploymentSafeguardResource struct{}

func TestAccKubernetesClusterDeploymentSafeguards_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_deployment_safeguard", "test")
	r := KubernetesClusterDeploymentSafeguardResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesClusterDeploymentSafeguards_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_deployment_safeguard", "test")
	r := KubernetesClusterDeploymentSafeguardResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccKubernetesClusterDeploymentSafeguards_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_deployment_safeguard", "test")
	r := KubernetesClusterDeploymentSafeguardResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesClusterDeploymentSafeguards_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_deployment_safeguard", "test")
	r := KubernetesClusterDeploymentSafeguardResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r KubernetesClusterDeploymentSafeguardResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	kubernetesClusterId, err := commonids.ParseKubernetesClusterID(state.ID)
	if err != nil {
		return nil, err
	}

	scopeId := commonids.NewScopeID(kubernetesClusterId.ID())

	resp, err := clients.Containers.DeploymentSafeguardsClient.Get(ctx, scopeId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("reading Deployment Safeguards for %s: %+v", kubernetesClusterId, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r KubernetesClusterDeploymentSafeguardResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_cluster_deployment_safeguard" "test" {
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  level                 = "Warn"
}
`, r.template(data))
}

func (r KubernetesClusterDeploymentSafeguardResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_cluster_deployment_safeguard" "import" {
  kubernetes_cluster_id = azurerm_kubernetes_cluster_deployment_safeguard.test.kubernetes_cluster_id
  level                 = azurerm_kubernetes_cluster_deployment_safeguard.test.level
}
`, r.basic(data))
}

func (r KubernetesClusterDeploymentSafeguardResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_cluster_deployment_safeguard" "test" {
  kubernetes_cluster_id        = azurerm_kubernetes_cluster.test.id
  level                        = "Enforce"
  excluded_namespaces          = ["my-app-namespace", "another-namespace"]
  pod_security_standards_level = "Restricted"
}
`, r.template(data))
}

func (r KubernetesClusterDeploymentSafeguardResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%[2]d"
  location = "%[1]s"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[2]d"

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_DS2_v2"
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type = "SystemAssigned"
  }

  azure_policy_enabled = true
}
`, data.Locations.Primary, data.RandomInteger)
}
