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

type ServiceConnectorContainerAppResource struct{}

func (r ServiceConnectorContainerAppResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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

func TestAccServiceConnectorContainerAppCosmosdb_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_connection", "test")
	r := ServiceConnectorContainerAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.cosmosdbBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccServiceConnectorContainerAppCosmosdb_secretAuth(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_connection", "test")
	r := ServiceConnectorContainerAppResource{}

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

func TestAccServiceConnectorContainerAppCosmosdb_servicePrincipalSecretAuth(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_connection", "test")
	r := ServiceConnectorContainerAppResource{}

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

func TestAccServiceConnectorContainerAppCosmosdb_userAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_connection", "test")
	r := ServiceConnectorContainerAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.cosmosdbWithUserAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("authentication"),
	})
}

func TestAccServiceConnectorContainerAppStorageBlob_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_connection", "test")
	r := ServiceConnectorContainerAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.storageBlob(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccServiceConnectorContainerAppStorageBlob_secretStore(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_connection", "test")
	r := ServiceConnectorContainerAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.secretStore(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccServiceConnectorContainerApp_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_connection", "test")
	r := ServiceConnectorContainerAppResource{}

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

func (r ServiceConnectorContainerAppResource) storageBlob(data acceptance.TestData) string {
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

resource "azurerm_container_app_environment" "test" {
  name                = "acctest-cae-%[3]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_container_app" "test" {
  name                         = "acctest-ca-%[3]d"
  container_app_environment_id = azurerm_container_app_environment.test.id
  resource_group_name          = azurerm_resource_group.test.name
  revision_mode                = "Single"

  template {
    container {
      name   = "acctest-cont-%[3]d"
      image  = "mcr.microsoft.com/k8se/quickstart:latest"
      cpu    = 0.25
      memory = "0.5Gi"
    }
  }

  lifecycle {
    ignore_changes = [
      template[0].container[0].env,
      identity,
    ]
  }
}

resource "azurerm_container_app_connection" "test" {
  name               = "acctestserviceconnector%[3]d"
  container_app_id   = azurerm_container_app.test.id
  target_resource_id = azurerm_storage_account.test.id
  authentication {
    type = "systemAssignedIdentity"
  }
}
`, data.Locations.Primary, data.RandomString, data.RandomInteger)
}

func (r ServiceConnectorContainerAppResource) cosmosdbBasic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_container_app_connection" "test" {
  name               = "acctestserviceconnector%[3]d"
  container_app_id   = azurerm_container_app.test.id
  target_resource_id = azurerm_cosmosdb_sql_database.test.id
  scope              = azurerm_container_app.test.template[0].container[0].name
  authentication {
    type = "systemAssignedIdentity"
  }
}
`, template, data.RandomString, data.RandomInteger)
}

func (r ServiceConnectorContainerAppResource) cosmosdbWithSecretAuth(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_container_app_connection" "test" {
  name               = "acctestserviceconnector%[2]d"
  container_app_id   = azurerm_container_app.test.id
  target_resource_id = azurerm_cosmosdb_sql_database.test.id
  scope              = azurerm_container_app.test.template[0].container[0].name
  authentication {
    type   = "secret"
    name   = "foo"
    secret = "bar"
  }
}
`, template, data.RandomInteger)
}

func (r ServiceConnectorContainerAppResource) cosmosdbWithServicePrincipalSecretAuth(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_container_app_connection" "test" {
  name               = "acctestserviceconnector%[3]d"
  container_app_id   = azurerm_container_app.test.id
  target_resource_id = azurerm_cosmosdb_sql_database.test.id
  scope              = azurerm_container_app.test.template[0].container[0].name
  authentication {
    type         = "servicePrincipalSecret"
    client_id    = "someclientid"
    principal_id = azurerm_user_assigned_identity.test.principal_id
    secret       = "somesecret"
  }
}
`, template, data.RandomString, data.RandomInteger)
}

func (r ServiceConnectorContainerAppResource) cosmosdbWithUserAssignedIdentity(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

data "azurerm_subscription" "test" {}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_container_app_connection" "test" {
  name               = "acctestserviceconnector%[3]d"
  container_app_id   = azurerm_container_app.test.id
  target_resource_id = azurerm_cosmosdb_sql_database.test.id
  scope              = azurerm_container_app.test.template[0].container[0].name
  authentication {
    type            = "userAssignedIdentity"
    subscription_id = data.azurerm_subscription.test.subscription_id
    client_id       = azurerm_user_assigned_identity.test.client_id
  }
}
`, template, data.RandomString, data.RandomInteger)
}

func (r ServiceConnectorContainerAppResource) secretStore(data acceptance.TestData) string {
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

resource "azurerm_container_app_environment" "test" {
  name                = "acctest-cae-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_container_app" "test" {
  name                         = "acctest-ca-%[3]d"
  container_app_environment_id = azurerm_container_app_environment.test.id
  resource_group_name          = azurerm_resource_group.test.name
  revision_mode                = "Single"

  template {
    container {
      name   = "acctest-cont-%[3]d"
      image  = "mcr.microsoft.com/k8se/quickstart:latest"
      cpu    = 0.25
      memory = "0.5Gi"
    }
  }

  lifecycle {
    ignore_changes = [
      template[0].container[0].env,
      identity,
      secret,
    ]
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

resource "azurerm_container_app_connection" "test" {
  name               = "acctestserviceconnector%[3]d"
  container_app_id   = azurerm_container_app.test.id
  target_resource_id = azurerm_storage_account.test.id
  client_type        = "java"

  secret_store {
    key_vault_id = azurerm_key_vault.test.id
  }
  authentication {
    type = "systemAssignedIdentity"
  }
}
`, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomString)
}

func (r ServiceConnectorContainerAppResource) complete(data acceptance.TestData) string {
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
  name                = "test-container%[4]s"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
  database_name       = azurerm_cosmosdb_sql_database.test.name
  partition_key_paths = ["/definition"]
}

resource "azurerm_container_app_environment" "test" {
  name                = "acctest-cae-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_virtual_network" "test" {
  name                = "vnet-%[3]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test1" {
  name                              = "subnet1"
  resource_group_name               = azurerm_resource_group.test.name
  virtual_network_name              = azurerm_virtual_network.test.name
  address_prefixes                  = ["10.0.1.0/24"]
  private_endpoint_network_policies = "Enabled"

  delegation {
    name = "delegation"

    service_delegation {
      name    = "Microsoft.App/environments"
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
    }
  }
}

resource "azurerm_container_app" "test" {
  name                         = "acctest-ca-%[3]d"
  container_app_environment_id = azurerm_container_app_environment.test.id
  resource_group_name          = azurerm_resource_group.test.name
  revision_mode                = "Single"

  template {
    container {
      name   = "acctest-cont-%[3]d"
      image  = "mcr.microsoft.com/k8se/quickstart:latest"
      cpu    = 0.25
      memory = "0.5Gi"
    }
  }

  lifecycle {
    ignore_changes = [
      template[0].container[0].env,
      identity,
      secret,
    ]
  }
}

resource "azurerm_container_app_connection" "test" {
  name               = "acctestserviceconnector%[3]d"
  container_app_id   = azurerm_container_app.test.id
  target_resource_id = azurerm_cosmosdb_sql_container.test.id
  client_type        = "springBoot"
  vnet_solution      = "serviceEndpoint"

  authentication {
    type = "systemAssignedIdentity"
  }
}
`, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomString)
}

func (r ServiceConnectorContainerAppResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[2]d"
  location = "%[1]s"
}

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctestacc%[3]s"
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

resource "azurerm_container_app_environment" "test" {
  name                = "acctest-cae-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_container_app" "test" {
  name                         = "acctest-ca-%[2]d"
  container_app_environment_id = azurerm_container_app_environment.test.id
  resource_group_name          = azurerm_resource_group.test.name
  revision_mode                = "Single"

  template {
    container {
      name   = "acctest-cont-%[2]d"
      image  = "mcr.microsoft.com/k8se/quickstart:latest"
      cpu    = 0.25
      memory = "0.5Gi"
    }
  }

  lifecycle {
    ignore_changes = [
      template[0].container[0].env,
      identity,
      secret,
    ]
  }
}
`, data.Locations.Primary, data.RandomInteger, data.RandomString)
}
