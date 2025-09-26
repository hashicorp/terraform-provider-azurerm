// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package serviceconnector_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicelinker/2024-04-01/servicelinker"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ServiceConnectorKubernetesClusterResource struct{}

func (r ServiceConnectorKubernetesClusterResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := servicelinker.ParseScopedLinkerID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.ServiceConnector.ServiceLinkerClient.LinkerGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return utils.Bool(true), nil
}

func TestAccServiceConnectorKubernetesClusterCosmosdb_secretAuth(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_connection", "test")
	r := ServiceConnectorKubernetesClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.cosmosdbWithSecretAuth(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("authentication"),
	})
}

func TestAccServiceConnectorKubernetesClusterCosmosdb_servicePrincipalSecretAuth(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_connection", "test")
	r := ServiceConnectorKubernetesClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.cosmosdbWithServicePrincipalSecretAuth(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("authentication"),
	})
}

func TestAccServiceConnectorKubernetesClusterStorageBlob_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_connection", "test")
	r := ServiceConnectorKubernetesClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.storageBlob(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("authentication"),
	})
}

func TestAccServiceConnectorKubernetesClusterStorageBlob_secretStore(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_connection", "test")
	r := ServiceConnectorKubernetesClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.secretStore(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("authentication"),
	})
}

func TestAccServiceConnectorKubernetesCluster_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_connection", "test")
	r := ServiceConnectorKubernetesClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("authentication"),
	})
}

func (r ServiceConnectorKubernetesClusterResource) storageBlob(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[3]d"
  location = "%[1]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestacc%[2]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%[3]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[3]d"

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "standard_a2_v2"

    upgrade_settings {
      max_surge                     = "10%%"
      drain_timeout_in_minutes      = 0
      node_soak_duration_in_minutes = 0
    }
  }

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_kubernetes_cluster_connection" "test" {
  name                  = "acctestserviceconnector%[3]d"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  target_resource_id    = azurerm_storage_account.test.id
  authentication {
    type   = "secret"
    name   = "foo"
    secret = "bar"
  }
}
`, data.Locations.Primary, data.RandomString, data.RandomInteger)
}

func (r ServiceConnectorKubernetesClusterResource) cosmosdbWithSecretAuth(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_kubernetes_cluster_connection" "test" {
  name                  = "acctestserviceconnector%[2]d"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  target_resource_id    = azurerm_cosmosdb_sql_database.test.id
  authentication {
    type   = "secret"
    name   = "foo"
    secret = "bar"
  }
}
`, template, data.RandomInteger)
}

func (r ServiceConnectorKubernetesClusterResource) cosmosdbWithServicePrincipalSecretAuth(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_kubernetes_cluster_connection" "test" {
  name                  = "acctestserviceconnector%[3]d"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  target_resource_id    = azurerm_cosmosdb_sql_database.test.id
  authentication {
    type         = "servicePrincipalSecret"
    client_id    = "someclientid"
    principal_id = azurerm_user_assigned_identity.test.principal_id
    secret       = "somesecret"
  }
}
`, template, data.RandomString, data.RandomInteger)
}

func (r ServiceConnectorKubernetesClusterResource) secretStore(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      recover_soft_deleted_key_vaults    = false
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[2]d"
  location = "%[1]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestacc%[4]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[2]d"

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "standard_a2_v2"

    upgrade_settings {
      max_surge                     = "10%%"
      drain_timeout_in_minutes      = 0
      node_soak_duration_in_minutes = 0
    }
  }

  identity {
    type = "SystemAssigned"
  }

  key_vault_secrets_provider {
    secret_rotation_enabled = true
  }
}

resource "azurerm_key_vault" "test" {
  name                     = "accAKV-%[4]s"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  tenant_id                = data.azurerm_client_config.current.tenant_id
  sku_name                 = "standard"
  purge_protection_enabled = true
}

resource "azurerm_kubernetes_cluster_connection" "test" {
  name                  = "acctestserviceconnector%[3]d"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  target_resource_id    = azurerm_storage_account.test.id
  client_type           = "java"

  secret_store {
    key_vault_id = azurerm_key_vault.test.id
  }
  authentication {
    type   = "secret"
    name   = "foo"
    secret = "bar"
  }
}
`, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomString)
}

func (r ServiceConnectorKubernetesClusterResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[2]d"
  location = "%[1]s"
}

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctestacc%[4]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "GlobalDocumentDB"

  consistency_policy {
    consistency_level       = "BoundedStaleness"
    max_interval_in_seconds = 10
    max_staleness_prefix    = 200
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}

resource "azurerm_cosmosdb_sql_database" "test" {
  name                = "cosmos-sql-db"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
  throughput          = 400
}

resource "azurerm_cosmosdb_sql_container" "test" {
  name                = "example-container"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
  database_name       = azurerm_cosmosdb_sql_database.test.name
  partition_key_paths = ["/definition"]
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[2]d"

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "standard_a2_v2"

    upgrade_settings {
      max_surge                     = "10%%"
      drain_timeout_in_minutes      = 0
      node_soak_duration_in_minutes = 0
    }
  }

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_kubernetes_cluster_connection" "test" {
  name                  = "acctestserviceconnector%[3]d"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  target_resource_id    = azurerm_cosmosdb_sql_database.test.id
  client_type           = "java"
  vnet_solution         = "privateLink"
  authentication {
    type   = "secret"
    name   = "foo"
    secret = "bar"
  }
}
`, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomString)
}

func (r ServiceConnectorKubernetesClusterResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[2]d"
  location = "%[1]s"
}

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctestacc%[4]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "GlobalDocumentDB"

  consistency_policy {
    consistency_level       = "BoundedStaleness"
    max_interval_in_seconds = 10
    max_staleness_prefix    = 200
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}

resource "azurerm_cosmosdb_sql_database" "test" {
  name                = "cosmos-sql-db"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
  throughput          = 400
}

resource "azurerm_cosmosdb_sql_container" "test" {
  name                = "example-container"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
  database_name       = azurerm_cosmosdb_sql_database.test.name
  partition_key_paths = ["/definition"]
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[2]d"

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "standard_a2_v2"

    upgrade_settings {
      max_surge                     = "10%%"
      drain_timeout_in_minutes      = 0
      node_soak_duration_in_minutes = 0
    }
  }

  identity {
    type = "SystemAssigned"
  }
}
`, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomString)
}
