package machinelearning_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/machinelearning/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type InferenceClusterResource struct{}

func TestAccInferenceCluster_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_inference_cluster", "test")
	r := InferenceClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.principal_id").Exists(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").Exists(),
			),
		},
		data.ImportStep("cluster_purpose", "description"),
	})
}

func TestAccInferenceCluster_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_inference_cluster", "test")
	r := InferenceClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.principal_id").Exists(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").Exists(),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccInferenceCluster_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_inference_cluster", "test")
	r := InferenceClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.principal_id").Exists(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").Exists(),
			),
		},
		data.ImportStep("ssl", "cluster_purpose", "description"),
	})
}

func TestAccInferenceCluster_basicUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_inference_cluster", "test")
	r := InferenceClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.principal_id").Exists(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").Exists(),
			),
		},
		data.ImportStep("cluster_purpose", "description"),
		{
			Config: r.basicUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.principal_id").Exists(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").Exists(),
			),
		},
		data.ImportStep("cluster_purpose", "description"),
	})
}

func TestAccInferenceCluster_completeUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_inference_cluster", "test")
	r := InferenceClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.principal_id").Exists(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").Exists(),
			),
		},
		data.ImportStep("ssl", "cluster_purpose", "description"),
		{
			Config: r.completeUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.principal_id").Exists(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").Exists(),
			),
		},
		data.ImportStep("ssl", "cluster_purpose", "description"),
	})
}

func (r InferenceClusterResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	inferenceClusterClient := client.MachineLearning.MachineLearningComputeClient
	id, err := parse.InferenceClusterID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := inferenceClusterClient.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.ComputeName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Inference Cluster %q: %+v", state.ID, err)
	}

	return utils.Bool(resp.Properties != nil), nil
}

func (r InferenceClusterResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_machine_learning_inference_cluster" "test" {
  name                          = "AIC-%d"
  machine_learning_workspace_id = azurerm_machine_learning_workspace.test.id
  location                      = azurerm_resource_group.test.location
  kubernetes_cluster_id         = azurerm_kubernetes_cluster.test.id
  cluster_purpose               = "Dev"
  node_pool_id                  = azurerm_kubernetes_cluster.test.default_node_pool[0].id

  identity {
    type = "SystemAssigned"
  }
}
`, template, data.RandomIntOfLength(8))
}

func (r InferenceClusterResource) basicUpdate(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_machine_learning_inference_cluster" "test" {
  name                          = "AIC-%d"
  machine_learning_workspace_id = azurerm_machine_learning_workspace.test.id
  location                      = azurerm_resource_group.test.location
  kubernetes_cluster_id         = azurerm_kubernetes_cluster.test.id
  cluster_purpose               = "Dev"
  node_pool_id                  = azurerm_kubernetes_cluster.test.default_node_pool[0].id

  identity {
    type = "SystemAssigned"
  }

  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomIntOfLength(8))
}

func (r InferenceClusterResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_machine_learning_inference_cluster" "test" {
  name                          = "AIC-%d"
  machine_learning_workspace_id = azurerm_machine_learning_workspace.test.id
  location                      = azurerm_resource_group.test.location
  kubernetes_cluster_id         = azurerm_kubernetes_cluster.test.id
  cluster_purpose               = "Test"
  node_pool_id                  = azurerm_kubernetes_cluster.test.default_node_pool[0].id
  ssl {
    cert  = file("testdata/cert.pem")
    key   = file("testdata/key.pem")
    cname = "www.contoso.com"
  }

  identity {
    type = "SystemAssigned"
  }

}
`, template, data.RandomIntOfLength(8))
}

func (r InferenceClusterResource) completeUpdate(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_machine_learning_inference_cluster" "test" {
  name                          = "AIC-%d"
  machine_learning_workspace_id = azurerm_machine_learning_workspace.test.id
  location                      = azurerm_resource_group.test.location
  kubernetes_cluster_id         = azurerm_kubernetes_cluster.test.id
  cluster_purpose               = "Test"
  node_pool_id                  = azurerm_kubernetes_cluster.test.default_node_pool[0].id
  ssl {
    cert  = file("testdata/cert.pem")
    key   = file("testdata/key.pem")
    cname = "www.contoso.com"
  }

  identity {
    type = "SystemAssigned"
  }

  tags = {
    ENV = "Test"
  }

}
`, template, data.RandomIntOfLength(8))
}

func (r InferenceClusterResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_machine_learning_inference_cluster" "import" {
  name                          = azurerm_machine_learning_inference_cluster.test.name
  machine_learning_workspace_id = azurerm_machine_learning_inference_cluster.test.machine_learning_workspace_id
  location                      = azurerm_machine_learning_inference_cluster.test.location
  kubernetes_cluster_id         = azurerm_machine_learning_inference_cluster.test.kubernetes_cluster_id
  node_pool_id                  = azurerm_machine_learning_inference_cluster.test.node_pool_id
  cluster_purpose               = azurerm_machine_learning_inference_cluster.test.cluster_purpose

  identity {
    type = "SystemAssigned"
  }

  tags = azurerm_machine_learning_inference_cluster.test.tags
}
`, template)
}

func (r InferenceClusterResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-ml-%[1]d"
  location = "%[2]s"
  tags = {
    "stage" = "test"
  }
}

resource "azurerm_application_insights" "test" {
  name                = "acctestai-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}

resource "azurerm_key_vault" "test" {
  name                = "acctestvault%[3]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id

  sku_name = "standard"

  purge_protection_enabled = true
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[4]d"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_machine_learning_workspace" "test" {
  name                    = "acctest-MLW%[5]d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  application_insights_id = azurerm_application_insights.test.id
  key_vault_id            = azurerm_key_vault.test.id
  storage_account_id      = azurerm_storage_account.test.id

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%[6]d"
  address_space       = ["10.1.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%[7]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.1.0.0/24"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = join("", ["acctestaks", azurerm_resource_group.test.location])
  node_resource_group = "acctestRGAKS-%d"

  default_node_pool {
    name           = "default"
    node_count     = 3
    vm_size        = "Standard_D11_v2"
    vnet_subnet_id = azurerm_subnet.test.id
  }

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary,
		data.RandomIntOfLength(12), data.RandomIntOfLength(15), data.RandomIntOfLength(16),
		data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
