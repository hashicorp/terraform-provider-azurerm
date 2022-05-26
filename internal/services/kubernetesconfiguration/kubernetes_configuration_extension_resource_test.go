package kubernetesconfiguration_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kubernetesconfiguration/sdk/2022-03-01/extensions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type KubernetesConfigurationExtensionResource struct{}

func TestAccKubernetesConfigurationExtension_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_configuration_extension", "test")
	r := KubernetesConfigurationExtensionResource{}
	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesConfigurationExtension_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_configuration_extension", "test")
	r := KubernetesConfigurationExtensionResource{}
	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccKubernetesConfigurationExtension_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_configuration_extension", "test")
	r := KubernetesConfigurationExtensionResource{}
	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("configuration_protected_settings"),
	})
}

func TestAccKubernetesConfigurationExtension_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_configuration_extension", "test")
	r := KubernetesConfigurationExtensionResource{}
	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
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

func (r KubernetesConfigurationExtensionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := extensions.ParseExtensionID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.KubernetesConfiguration.ExtensionsClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(true), nil
}

func (r KubernetesConfigurationExtensionResource) template(data acceptance.TestData) string {
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

func (r KubernetesConfigurationExtensionResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_kubernetes_configuration_extension" "test" {
  name                  = "acctest-kc-%d"
  resource_group_name   = azurerm_resource_group.test.name
  cluster_name          = azurerm_kubernetes_cluster.test.name
  cluster_resource_name = "managedClusters"
  extension_type        = "microsoft.flux"
}
`, template, data.RandomInteger)
}

func (r KubernetesConfigurationExtensionResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_kubernetes_configuration_extension" "import" {
  name                  = azurerm_kubernetes_configuration_extension.test.name
  resource_group_name   = azurerm_kubernetes_configuration_extension.test.resource_group_name
  cluster_name          = azurerm_kubernetes_configuration_extension.test.cluster_name
  cluster_resource_name = azurerm_kubernetes_configuration_extension.test.cluster_resource_name
  extension_type        = azurerm_kubernetes_configuration_extension.test.extension_type
}
`, config)
}

func (r KubernetesConfigurationExtensionResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_kubernetes_configuration_extension" "test" {
  name                  = "acctest-kc-%d"
  resource_group_name   = azurerm_resource_group.test.name
  cluster_name          = azurerm_kubernetes_cluster.test.name
  cluster_resource_name = "managedClusters"
  extension_type        = "microsoft.flux"
  version               = "1.2.0"
  release_namespace     = "release1"

  configuration_protected_settings = {
    "omsagent.secret.key" = "secretKeyValue01"
  }

  configuration_settings = {
    "omsagent.secret.wsid"     = "a38cef99-5a89-52ed-b6db-22095c23664b",
    "omsagent.env.clusterName" = "clusterName1"
  }
}
`, template, data.RandomInteger)
}

func (r KubernetesConfigurationExtensionResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_kubernetes_configuration_extension" "test" {
  name                  = "acctest-kc-%d"
  resource_group_name   = azurerm_resource_group.test.name
  cluster_name          = azurerm_kubernetes_cluster.test.name
  cluster_resource_name = "managedClusters"
  extension_type        = "microsoft.flux"
  version               = "1.2.0"
  release_namespace     = "release1"

  configuration_protected_settings = {
    "omsagent.secret.key" = "secretKeyValue02"
  }

  configuration_settings = {
    "omsagent.secret.wsid"     = "a38cef99-5a89-52ed-b6db-22095c23664c",
    "omsagent.env.clusterName" = "clusterName2"
  }
}
`, template, data.RandomInteger)
}
