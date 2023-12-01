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

type ServiceConnectorAppServiceResource struct{}

func (r ServiceConnectorAppServiceResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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

func TestAccServiceConnectorAppServiceCosmosdb_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_connection", "test")
	r := ServiceConnectorAppServiceResource{}

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

func TestAccServiceConnectorAppServiceCosmosdb_secretAuth(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_connection", "test")
	r := ServiceConnectorAppServiceResource{}

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

func TestAccServiceConnectorAppServiceCosmosdb_servicePrincipalSecretAuth(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_connection", "test")
	r := ServiceConnectorAppServiceResource{}

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

func TestAccServiceConnectorAppServiceCosmosdb_userAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_connection", "test")
	r := ServiceConnectorAppServiceResource{}

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

func TestAccServiceConnectorAppServiceStorageBlob_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_connection", "test")
	r := ServiceConnectorAppServiceResource{}

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

func TestAccServiceConnectorAppServiceStorageBlob_secretStore(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_connection", "test")
	r := ServiceConnectorAppServiceResource{}

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

func TestAccServiceConnectorAppService_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_connection", "test")
	r := ServiceConnectorAppServiceResource{}

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

func (r ServiceConnectorAppServiceResource) storageBlob(data acceptance.TestData) string {
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

resource "azurerm_service_plan" "test" {
  location            = azurerm_resource_group.test.location
  name                = "testserviceplan%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "P1v2"
  os_type             = "Linux"
}

resource "azurerm_linux_web_app" "test" {
  location            = azurerm_resource_group.test.location
  name                = "testlinuxwebapp%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  site_config {}

  lifecycle {
    ignore_changes = [
      app_settings["AZURE_STORAGEBLOB_RESOURCEENDPOINT"],
      identity,
      sticky_settings,
    ]
  }
}

resource "azurerm_app_service_connection" "test" {
  name               = "acctestserviceconnector%[3]d"
  app_service_id     = azurerm_linux_web_app.test.id
  target_resource_id = azurerm_storage_account.test.id
  authentication {
    type = "systemAssignedIdentity"
  }
}
`, data.Locations.Primary, data.RandomString, data.RandomInteger)
}

func (r ServiceConnectorAppServiceResource) cosmosdbBasic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_app_service_connection" "test" {
  name               = "acctestserviceconnector%[3]d"
  app_service_id     = azurerm_linux_web_app.test.id
  target_resource_id = azurerm_cosmosdb_sql_database.test.id
  authentication {
    type = "systemAssignedIdentity"
  }
}
`, template, data.RandomString, data.RandomInteger)
}

func (r ServiceConnectorAppServiceResource) cosmosdbWithSecretAuth(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_app_service_connection" "test" {
  name               = "acctestserviceconnector%[2]d"
  app_service_id     = azurerm_linux_web_app.test.id
  target_resource_id = azurerm_cosmosdb_sql_database.test.id
  authentication {
    type   = "secret"
    name   = "foo"
    secret = "bar"
  }
}
`, template, data.RandomInteger)
}

func (r ServiceConnectorAppServiceResource) cosmosdbWithServicePrincipalSecretAuth(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_app_service_connection" "test" {
  name               = "acctestserviceconnector%[3]d"
  app_service_id     = azurerm_linux_web_app.test.id
  target_resource_id = azurerm_cosmosdb_sql_database.test.id
  authentication {
    type         = "servicePrincipalSecret"
    client_id    = "someclientid"
    principal_id = azurerm_user_assigned_identity.test.principal_id
    secret       = "somesecret"
  }
}
`, template, data.RandomString, data.RandomInteger)
}

func (r ServiceConnectorAppServiceResource) cosmosdbWithUserAssignedIdentity(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

data "azurerm_subscription" "test" {}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_app_service_connection" "test" {
  name               = "acctestserviceconnector%[3]d"
  app_service_id     = azurerm_linux_web_app.test.id
  target_resource_id = azurerm_cosmosdb_sql_database.test.id
  authentication {
    type            = "userAssignedIdentity"
    subscription_id = data.azurerm_subscription.test.subscription_id
    client_id       = azurerm_user_assigned_identity.test.client_id
  }
}
`, template, data.RandomString, data.RandomInteger)
}

func (r ServiceConnectorAppServiceResource) secretStore(data acceptance.TestData) string {
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

resource "azurerm_service_plan" "test" {
  name                = "acctestASP-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  os_type             = "Linux"
  sku_name            = "P1v2"
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
      name    = "Microsoft.Web/serverFarms"
      actions = ["Microsoft.Network/virtualNetworks/subnets/action"]
    }
  }
}


resource "azurerm_linux_web_app" "test" {
  name                      = "acctestWA-%[3]d"
  location                  = azurerm_resource_group.test.location
  resource_group_name       = azurerm_resource_group.test.name
  service_plan_id           = azurerm_service_plan.test.id
  virtual_network_subnet_id = azurerm_subnet.test1.id

  site_config {}
  lifecycle {
    ignore_changes = [
      app_settings["AZURE_STORAGEBLOB_RESOURCEENDPOINT"],
      identity,
      sticky_settings,
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

resource "azurerm_app_service_connection" "test" {
  name               = "acctestserviceconnector%[3]d"
  app_service_id     = azurerm_linux_web_app.test.id
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

func (r ServiceConnectorAppServiceResource) complete(data acceptance.TestData) string {
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

resource "azurerm_service_plan" "test" {
  name                = "acctestASP-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  os_type             = "Linux"
  sku_name            = "B1"
}

resource "azurerm_virtual_network" "test" {
  name                = "vnet-%[3]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test1" {
  name                 = "subnet1"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]

  delegation {
    name = "delegation"

    service_delegation {
      name    = "Microsoft.Web/serverFarms"
      actions = ["Microsoft.Network/virtualNetworks/subnets/action"]
    }
  }

  lifecycle {
    ignore_changes = [
      service_endpoints,
    ]
  }
}

resource "azurerm_linux_web_app" "test" {
  name                      = "acctestWA-%[3]d"
  location                  = azurerm_resource_group.test.location
  resource_group_name       = azurerm_resource_group.test.name
  service_plan_id           = azurerm_service_plan.test.id
  virtual_network_subnet_id = azurerm_subnet.test1.id

  site_config {}

  lifecycle {
    ignore_changes = [
      app_settings,
      identity,
      sticky_settings,
    ]
  }
}

resource "azurerm_app_service_connection" "test" {
  name               = "acctestserviceconnector%[3]d"
  app_service_id     = azurerm_linux_web_app.test.id
  target_resource_id = azurerm_cosmosdb_sql_database.test.id
  client_type        = "java"
  vnet_solution      = "serviceEndpoint"
  authentication {
    type = "systemAssignedIdentity"
  }
}
`, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomString)
}

func (r ServiceConnectorAppServiceResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
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

resource "azurerm_cosmosdb_sql_container" "test" {
  name                = "test-container%[3]s"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
  database_name       = azurerm_cosmosdb_sql_database.test.name
  partition_key_paths = ["/definition"]
}

resource "azurerm_service_plan" "test" {
  location            = azurerm_resource_group.test.location
  name                = "testserviceplan%[3]s"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "P1v2"
  os_type             = "Linux"
}

resource "azurerm_linux_web_app" "test" {
  location            = azurerm_resource_group.test.location
  name                = "linuxwebapp%[3]s"
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id
  site_config {}

  lifecycle {
    ignore_changes = [
      app_settings,
      identity,
      sticky_settings,
    ]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
