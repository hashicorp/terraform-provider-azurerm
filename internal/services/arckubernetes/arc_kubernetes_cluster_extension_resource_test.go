// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package arckubernetes_test

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

type ArcKubernetesClusterExtensionResource struct{}

func TestAccArcKubernetesClusterExtension_basic(t *testing.T) {
	credential, privateKey, publicKey := ArcKubernetesClusterResource{}.getCredentials(t)
	data := acceptance.BuildTestData(t, "azurerm_arc_kubernetes_cluster_extension", "test")
	r := ArcKubernetesClusterExtensionResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, credential, privateKey, publicKey),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccArcKubernetesClusterExtension_requiresImport(t *testing.T) {
	credential, privateKey, publicKey := ArcKubernetesClusterResource{}.getCredentials(t)
	data := acceptance.BuildTestData(t, "azurerm_arc_kubernetes_cluster_extension", "test")
	r := ArcKubernetesClusterExtensionResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, credential, privateKey, publicKey),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data, credential, privateKey, publicKey),
			ExpectError: acceptance.RequiresImportError(data.ResourceType),
		},
	})
}

func TestAccArcKubernetesClusterExtension_complete(t *testing.T) {
	credential, privateKey, publicKey := ArcKubernetesClusterResource{}.getCredentials(t)
	data := acceptance.BuildTestData(t, "azurerm_arc_kubernetes_cluster_extension", "test")
	r := ArcKubernetesClusterExtensionResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, credential, privateKey, publicKey),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("configuration_protected_settings"),
	})
}

func TestAccArcKubernetesClusterExtension_update(t *testing.T) {
	credential, privateKey, publicKey := ArcKubernetesClusterResource{}.getCredentials(t)
	data := acceptance.BuildTestData(t, "azurerm_arc_kubernetes_cluster_extension", "test")
	r := ArcKubernetesClusterExtensionResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, credential, privateKey, publicKey),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("configuration_protected_settings"),
		{
			Config: r.update(data, credential, privateKey, publicKey),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("configuration_protected_settings"),
	})
}

func (r ArcKubernetesClusterExtensionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := extensions.ParseScopedExtensionID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.ArcKubernetes.ExtensionsClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r ArcKubernetesClusterExtensionResource) template(data acceptance.TestData, credential string, privateKey string, publicKey string) string {
	return fmt.Sprintf(`
				%[1]s

resource "azurerm_arc_kubernetes_cluster" "test" {
  name                         = "acctest-akcc-%[2]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  agent_public_key_certificate = "%[4]s"
  identity {
    type = "SystemAssigned"
  }

  %[3]s

  depends_on = [
    azurerm_linux_virtual_machine.test
  ]
}
`, ArcKubernetesClusterResource{}.template(data, credential), data.RandomInteger, ArcKubernetesClusterResource{}.provisionTemplate(data, credential, privateKey), publicKey)
}

func (r ArcKubernetesClusterExtensionResource) basic(data acceptance.TestData, credential string, privateKey string, publicKey string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_arc_kubernetes_cluster_extension" "test" {
  name           = "acctest-kce-%[2]d"
  cluster_id     = azurerm_arc_kubernetes_cluster.test.id
  extension_type = "microsoft.flux"

  identity {
    type = "SystemAssigned"
  }

  depends_on = [
    azurerm_linux_virtual_machine.test
  ]
}
`, r.template(data, credential, privateKey, publicKey), data.RandomInteger)
}

func (r ArcKubernetesClusterExtensionResource) requiresImport(data acceptance.TestData, credential string, privateKey string, publicKey string) string {
	config := r.basic(data, credential, privateKey, publicKey)
	return fmt.Sprintf(`
			%s

resource "azurerm_arc_kubernetes_cluster_extension" "import" {
  name           = azurerm_arc_kubernetes_cluster_extension.test.name
  cluster_id     = azurerm_arc_kubernetes_cluster_extension.test.cluster_id
  extension_type = azurerm_arc_kubernetes_cluster_extension.test.extension_type

  identity {
    type = "SystemAssigned"
  }

  depends_on = [
    azurerm_linux_virtual_machine.test
  ]
}
`, config)
}

func (r ArcKubernetesClusterExtensionResource) complete(data acceptance.TestData, credential string, privateKey string, publicKey string) string {
	return fmt.Sprintf(`
			%[1]s

resource "azurerm_arc_kubernetes_cluster_extension" "test" {
  name              = "acctest-kce-%[2]d"
  cluster_id        = azurerm_arc_kubernetes_cluster.test.id
  extension_type    = "microsoft.flux"
  version           = "1.6.3"
  release_namespace = "flux-system"

  configuration_protected_settings = {
    "omsagent.secret.key" = "secretKeyValue1"
  }

  configuration_settings = {
    "omsagent.env.clusterName" = "clusterName1"
  }

  identity {
    type = "SystemAssigned"
  }

  depends_on = [
    azurerm_linux_virtual_machine.test
  ]
}
`, r.template(data, credential, privateKey, publicKey), data.RandomInteger)
}

func (r ArcKubernetesClusterExtensionResource) update(data acceptance.TestData, credential string, privateKey string, publicKey string) string {
	return fmt.Sprintf(`
			%[1]s

resource "azurerm_arc_kubernetes_cluster_extension" "test" {
  name              = "acctest-kce-%[2]d"
  cluster_id        = azurerm_arc_kubernetes_cluster.test.id
  extension_type    = "microsoft.flux"
  version           = "1.6.3"
  release_namespace = "flux-system"

  configuration_protected_settings = {
    "omsagent.secret.key" = "secretKeyValue2"
  }

  configuration_settings = {
    "omsagent.env.clusterName" = "clusterName2"
  }

  identity {
    type = "SystemAssigned"
  }

  depends_on = [
    azurerm_linux_virtual_machine.test
  ]
}
`, r.template(data, credential, privateKey, publicKey), data.RandomInteger)
}
