package machinelearning_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/machinelearning/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type ComputeClusterResource struct{}

func TestAccComputeCluster_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_compute_cluster", "test")
	r := ComputeClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.principal_id").Exists(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").Exists(),
				check.That(data.ResourceName).Key("scale_settings.#").HasValue("1"),
				check.That(data.ResourceName).Key("scale_settings.0.max_node_count").Exists(),
				check.That(data.ResourceName).Key("scale_settings.0.min_node_count").Exists(),
				check.That(data.ResourceName).Key("scale_settings.0.scale_down_nodes_after_idle_duration").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccComputeCluster_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_compute_cluster", "test")
	r := ComputeClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.principal_id").Exists(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").Exists(),
				check.That(data.ResourceName).Key("scale_settings.#").HasValue("1"),
				check.That(data.ResourceName).Key("scale_settings.0.max_node_count").Exists(),
				check.That(data.ResourceName).Key("scale_settings.0.min_node_count").Exists(),
				check.That(data.ResourceName).Key("scale_settings.0.scale_down_nodes_after_idle_duration").Exists(),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (r ComputeClusterResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	computeClusterClient := client.MachineLearning.MachineLearningComputeClient
	id, err := parse.ComputeClusterID(state.ID)

	if err != nil {
		return nil, err
	}

	computeResource, err := computeClusterClient.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.ComputeName)
	if err != nil {
		if utils.ResponseWasNotFound(computeResource.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Machine Learning Compute Cluster %q: %+v", state.ID, err)
	}
	return utils.Bool(computeResource.Properties != nil), nil
}

func (r ComputeClusterResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_machine_learning_compute_cluster" "test" {
  name                          = "CC-%d"
  location                      = azurerm_resource_group.test.location
  vm_priority                   = "LowPriority"
  vm_size                       = "STANDARD_DS2_V2"
  machine_learning_workspace_id = azurerm_machine_learning_workspace.test.id
  subnet_resource_id            = azurerm_subnet.test.id

  scale_settings {
    min_node_count                       = 0
    max_node_count                       = 1
    scale_down_nodes_after_idle_duration = "PT30S" # 30 seconds
  }

  identity {
    type = "SystemAssigned"
  }
}
`, template, data.RandomIntOfLength(8))
}

func (r ComputeClusterResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_machine_learning_compute_cluster" "import" {
  name                          = azurerm_machine_learning_compute_cluster.test.name
  location                      = azurerm_machine_learning_compute_cluster.test.location
  vm_priority                   = azurerm_machine_learning_compute_cluster.test.vm_priority
  vm_size                       = azurerm_machine_learning_compute_cluster.test.vm_size
  machine_learning_workspace_id = azurerm_machine_learning_compute_cluster.test.machine_learning_workspace_id

  scale_settings {
    min_node_count                       = 0
    max_node_count                       = 1
    scale_down_nodes_after_idle_duration = "PT2M" # 120 seconds
  }

  identity {
    type = "SystemAssigned"
  }
}

`, template)
}

func (r ComputeClusterResource) template(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary,
		data.RandomIntOfLength(12), data.RandomIntOfLength(15), data.RandomIntOfLength(16),
		data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
