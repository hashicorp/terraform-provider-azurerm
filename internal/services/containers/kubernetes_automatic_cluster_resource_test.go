// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package containers_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type KubernetesAutomaticClusterResource struct{}

func TestAccKubernetesAutomaticCluster_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

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

func TestAccKubernetesAutomaticCluster_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImportConfig(data),
			ExpectError: acceptance.RequiresImportError("azurerm_kubernetes_automatic_cluster"),
		},
	})
}

func TestAccKubernetesAutomaticCluster_tags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.tagsConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.tagsUpdatedConfig(data),
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

	resp, err := clients.Containers.KubernetesClustersClient_v2026_04_01.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading Kubernetes Cluster (%s): %+v", id.String(), err)
	}

	return pointer.To(resp.Model != nil && resp.Model.Id != nil), nil
}

func (KubernetesAutomaticClusterResource) basic(data acceptance.TestData) string {
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

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r KubernetesAutomaticClusterResource) requiresImportConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_automatic_cluster" "import" {
  name                = azurerm_kubernetes_automatic_cluster.test.name
  location            = azurerm_kubernetes_automatic_cluster.test.location
  resource_group_name = azurerm_kubernetes_automatic_cluster.test.resource_group_name

  identity {
    type = "SystemAssigned"
  }
}
`, r.basic(data))
}

func (KubernetesAutomaticClusterResource) tagsConfig(data acceptance.TestData) string {
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

  identity {
    type = "SystemAssigned"
  }

  tags = {
    dimension = "C-137"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (KubernetesAutomaticClusterResource) tagsUpdatedConfig(data acceptance.TestData) string {
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

  identity {
    type = "SystemAssigned"
  }

  tags = {
    dimension = "D-99"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
