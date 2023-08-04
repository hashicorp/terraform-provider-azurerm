// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package machinelearning_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2023-04-01/machinelearningcomputes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SynapseSparkResource struct{}

func TestAccSynapseSpark_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_synapse_spark", "test")
	r := SynapseSparkResource{}

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

func TestAccSynapseSpark_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_synapse_spark", "test")
	r := SynapseSparkResource{}

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

func TestAccSynapseSpark_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_synapse_spark", "test")
	r := SynapseSparkResource{}

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

func TestAccSynapseSpark_identity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_synapse_spark", "test")
	r := SynapseSparkResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.identitySystemAssignedUserAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.0.principal_id").IsUUID(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.identityUserAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.identitySystemAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.0.principal_id").IsUUID(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func (r SynapseSparkResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	computeClient := client.MachineLearning.MachineLearningComputes
	id, err := machinelearningcomputes.ParseComputeID(state.ID)
	if err != nil {
		return nil, err
	}

	computeResource, err := computeClient.ComputeGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(computeResource.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Machine Learning Compute Cluster %q: %+v", state.ID, err)
	}
	return utils.Bool(computeResource.Model.Properties != nil), nil
}

func (r SynapseSparkResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_machine_learning_synapse_spark" "test" {
  name                          = "acctest%d"
  machine_learning_workspace_id = azurerm_machine_learning_workspace.test.id
  location                      = azurerm_resource_group.test.location
  synapse_spark_pool_id         = azurerm_synapse_spark_pool.test.id
  local_auth_enabled            = false
}
`, template, data.RandomIntOfLength(8))
}

func (r SynapseSparkResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s
resource "azurerm_machine_learning_synapse_spark" "test" {
  name                          = "acctest%d"
  machine_learning_workspace_id = azurerm_machine_learning_workspace.test.id
  location                      = azurerm_resource_group.test.location
  synapse_spark_pool_id         = azurerm_synapse_spark_pool.test.id
  description                   = "this is desc changed"
  identity {
    type = "SystemAssigned"
  }

  tags = {
    Ke1y = "Val1ue"
  }
}
`, template, data.RandomIntOfLength(8))
}

func (r SynapseSparkResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s
resource "azurerm_machine_learning_synapse_spark" "import" {
  name                          = azurerm_machine_learning_synapse_spark.test.name
  location                      = azurerm_machine_learning_synapse_spark.test.location
  machine_learning_workspace_id = azurerm_machine_learning_synapse_spark.test.machine_learning_workspace_id
  synapse_spark_pool_id         = azurerm_machine_learning_synapse_spark.test.synapse_spark_pool_id
}
`, template)
}

func (r SynapseSparkResource) identitySystemAssigned(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s
resource "azurerm_machine_learning_synapse_spark" "test" {
  name                          = "acctest%d"
  machine_learning_workspace_id = azurerm_machine_learning_workspace.test.id
  location                      = azurerm_resource_group.test.location
  synapse_spark_pool_id         = azurerm_synapse_spark_pool.test.id

  identity {
    type = "SystemAssigned"
  }
}
`, template, data.RandomIntOfLength(8))
}

func (r SynapseSparkResource) identityUserAssigned(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s
resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestUAI-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_machine_learning_synapse_spark" "test" {
  name                          = "acctest%d"
  machine_learning_workspace_id = azurerm_machine_learning_workspace.test.id
  location                      = azurerm_resource_group.test.location
  synapse_spark_pool_id         = azurerm_synapse_spark_pool.test.id

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }
}
`, template, data.RandomInteger, data.RandomIntOfLength(8))
}

func (r SynapseSparkResource) identitySystemAssignedUserAssigned(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s
resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestUAI-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_machine_learning_synapse_spark" "test" {
  name                          = "acctest%d"
  machine_learning_workspace_id = azurerm_machine_learning_workspace.test.id
  location                      = azurerm_resource_group.test.location
  synapse_spark_pool_id         = azurerm_synapse_spark_pool.test.id

  identity {
    type = "SystemAssigned, UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }
}
`, template, data.RandomInteger, data.RandomIntOfLength(8))
}

func (r SynapseSparkResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
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
  name                     = "acckv%[3]d"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  tenant_id                = data.azurerm_client_config.current.tenant_id
  sku_name                 = "standard"
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

resource "azurerm_storage_data_lake_gen2_filesystem" "test" {
  name               = "acc%d"
  storage_account_id = azurerm_storage_account.test.id
}

resource "azurerm_synapse_workspace" "test" {
  name                                 = "acctestsw%d"
  resource_group_name                  = azurerm_resource_group.test.name
  location                             = azurerm_resource_group.test.location
  storage_data_lake_gen2_filesystem_id = azurerm_storage_data_lake_gen2_filesystem.test.id
  sql_administrator_login              = "sqladminuser"
  sql_administrator_login_password     = "H@Sh1CoR3!"

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_synapse_spark_pool" "test" {
  name                 = "acc%d"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  node_size_family     = "MemoryOptimized"
  node_size            = "Small"
  node_count           = 3
}
`, data.RandomInteger, data.Locations.Primary,
		data.RandomInteger, data.RandomIntOfLength(15), data.RandomIntOfLength(16),
		data.RandomIntOfLength(8), data.RandomIntOfLength(8), data.RandomIntOfLength(8), data.RandomIntOfLength(8))
}
