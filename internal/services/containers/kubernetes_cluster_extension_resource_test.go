// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containers_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kubernetesconfiguration/2022-11-01/extensions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type KubernetesClusterExtensionResource struct{}

func TestAccKubernetesClusterExtension_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_extension", "test")
	r := KubernetesClusterExtensionResource{}
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

func TestAccKubernetesClusterExtension_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_extension", "test")
	r := KubernetesClusterExtensionResource{}
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

func TestAccKubernetesClusterExtension_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_extension", "test")
	r := KubernetesClusterExtensionResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("configuration_protected_settings"),
	})
}

func TestAccKubernetesClusterExtension_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_extension", "test")
	r := KubernetesClusterExtensionResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("configuration_protected_settings"),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("configuration_protected_settings"),
	})
}

func TestAccKubernetesClusterExtension_plan(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_extension", "test")
	r := KubernetesClusterExtensionResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.plan(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r KubernetesClusterExtensionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := extensions.ParseScopedExtensionID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Containers.KubernetesExtensionsClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r KubernetesClusterExtensionResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestAKC-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[1]d"

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_DS2_v2"
  }

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r KubernetesClusterExtensionResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_kubernetes_cluster_extension" "test" {
  name           = "acctest-kce-%d"
  cluster_id     = azurerm_kubernetes_cluster.test.id
  extension_type = "microsoft.flux"
}
`, template, data.RandomInteger)
}

func (r KubernetesClusterExtensionResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_kubernetes_cluster_extension" "import" {
  name           = azurerm_kubernetes_cluster_extension.test.name
  cluster_id     = azurerm_kubernetes_cluster_extension.test.cluster_id
  extension_type = azurerm_kubernetes_cluster_extension.test.extension_type
}
`, config)
}

func (r KubernetesClusterExtensionResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_kubernetes_cluster_extension" "test" {
  name              = "acctest-kce-%d"
  cluster_id        = azurerm_kubernetes_cluster.test.id
  extension_type    = "microsoft.flux"
  version           = "1.6.3"
  release_namespace = "flux-system"

  configuration_protected_settings = {
    "omsagent.secret.key" = "secretKeyValue1"
  }

  configuration_settings = {
    "omsagent.env.clusterName" = "clusterName1"
  }
}
`, template, data.RandomInteger)
}

func (r KubernetesClusterExtensionResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_kubernetes_cluster_extension" "test" {
  name              = "acctest-kce-%d"
  cluster_id        = azurerm_kubernetes_cluster.test.id
  extension_type    = "microsoft.flux"
  version           = "1.6.3"
  release_namespace = "flux-system"

  configuration_protected_settings = {
    "omsagent.secret.key" = "secretKeyValue2"
  }

  configuration_settings = {
    "omsagent.env.clusterName" = "clusterName2"
  }
}
`, template, data.RandomInteger)
}

func (r KubernetesClusterExtensionResource) plan(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_kubernetes_cluster_extension" "test" {
  name           = "acctest-kce-%d"
  cluster_id     = azurerm_kubernetes_cluster.test.id
  extension_type = "cognosys.nodejs-on-alpine"

  configuration_settings = {
    "title" = "Title",
  }

  plan {
    name      = "nodejs-18-alpine-container"
    product   = "nodejs18-alpine-container"
    publisher = "cognosys"
  }
}
`, template, data.RandomInteger)
}
