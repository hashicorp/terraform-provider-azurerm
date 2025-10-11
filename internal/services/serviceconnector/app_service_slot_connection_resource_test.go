// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package serviceconnector_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/servicelinker/2024-04-01/servicelinker"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type AppServiceSlotConnectorResource struct{}

func TestAccAppServiceSlotConnectorCosmosdb_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot_connection", "test")
	r := AppServiceSlotConnectorResource{}

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

func TestAccAppServiceSlotConnectorCosmosdb_secretAuth(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot_connection", "test")
	r := AppServiceSlotConnectorResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.cosmosdbSecretAuth(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("authentication"),
	})
}

func TestAccAppServiceSlotConnectorCosmosdb_servicePrincipalSecretAuth(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot_connection", "test")
	r := AppServiceSlotConnectorResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.cosmosdbServicePrincipalSecretAuth(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("authentication"),
	})
}

func TestAccAppServiceSlotConnectorCosmosdb_userAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot_connection", "test")
	r := AppServiceSlotConnectorResource{}

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

func TestAccAppServiceSlotConnectorStorageBlob_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot_connection", "test")
	r := AppServiceSlotConnectorResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.storageBlob(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("authentication.0.secret"),
	})
}

func TestAccAppServiceSlotConnectorStorageBlob_secretStore(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot_connection", "test")
	r := AppServiceSlotConnectorResource{}

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

func TestAccAppServiceSlotConnector_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot_connection", "test")
	r := AppServiceSlotConnectorResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.storageBlob(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccAppServiceSlotConnector_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot_connection", "test")
	r := AppServiceSlotConnectorResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("authentication.0.secret"),
	})
}

func (r AppServiceSlotConnectorResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := servicelinker.ParseScopedLinkerID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.ServiceConnector.ServiceLinkerClient.LinkerGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r AppServiceSlotConnectorResource) storageBlob(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  app_service_name    = azurerm_app_service.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  lifecycle {
    ignore_changes = [
      identity,
    ]
  }
}

resource "azurerm_app_service_slot_connection" "test" {
  name                = "acctestserviceconnector%d"
  app_service_slot_id = azurerm_app_service_slot.test.id
  target_resource_id  = azurerm_storage_account.test.id
  authentication {
    type = "systemAssignedIdentity"
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r AppServiceSlotConnectorResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_slot_connection" "import" {
  name                = azurerm_app_service_slot_connection.test.name
  app_service_slot_id = azurerm_app_service_slot_connection.test.app_service_slot_id
  target_resource_id  = azurerm_app_service_slot_connection.test.target_resource_id
  authentication {
    type = "systemAssignedIdentity"
  }
}
`, r.storageBlob(data))
}

func (r AppServiceSlotConnectorResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  app_service_name    = azurerm_app_service.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
}

resource "azurerm_app_service_slot_connection" "test" {
  name                = "acctestserviceconnector%d"
  app_service_slot_id = azurerm_app_service_slot.test.id
  target_resource_id  = azurerm_storage_account.test.id
  client_type         = "dotnet"
  authentication {
    type   = "secret"
    name   = "name"
    secret = "secret"
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r AppServiceSlotConnectorResource) storageBlobUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  app_service_name    = azurerm_app_service.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  lifecycle {
    ignore_changes = [
      identity,
    ]
  }
}

resource "azurerm_app_service_slot_connection" "test" {
  name                = "acctestserviceconnector%d"
  app_service_slot_id = azurerm_app_service_slot.test.id
  target_resource_id  = azurerm_storage_account.test.id
  client_type         = "python"
  authentication {
    type = "systemAssignedIdentity"
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r AppServiceSlotConnectorResource) cosmosdbBasic(data acceptance.TestData) string {
	template := r.templateWithCosmosDB(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%[2]d"
  app_service_name    = azurerm_app_service.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  lifecycle {
    ignore_changes = [
      identity,
    ]
  }
}

resource "azurerm_app_service_slot_connection" "test" {
  name                = "acctestserviceconnector%[2]d"
  app_service_slot_id = azurerm_app_service_slot.test.id
  target_resource_id  = azurerm_cosmosdb_account.test.id
  authentication {
    type = "systemAssignedIdentity"
  }
}
`, template, data.RandomInteger)
}

func (r AppServiceSlotConnectorResource) cosmosdbSecretAuth(data acceptance.TestData) string {
	template := r.templateWithCosmosDB(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%[2]d"
  app_service_name    = azurerm_app_service.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  lifecycle {
    ignore_changes = [
      identity,
    ]
  }
}

resource "azurerm_app_service_slot_connection" "test" {
  name                = "acctestserviceconnector%[2]d"
  app_service_slot_id = azurerm_app_service_slot.test.id
  target_resource_id  = azurerm_cosmosdb_account.test.id
  authentication {
    type   = "secret"
    name   = "foo"
    secret = "bar"
  }
}
`, template, data.RandomInteger)
}

func (r AppServiceSlotConnectorResource) cosmosdbServicePrincipalSecretAuth(data acceptance.TestData) string {
	template := r.templateWithCosmosDB(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%[3]d"
  app_service_name    = azurerm_app_service.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  lifecycle {
    ignore_changes = [
      identity,
    ]
  }
}

resource "azurerm_app_service_slot_connection" "test" {
  name                = "acctestserviceconnector%[3]d"
  app_service_slot_id = azurerm_app_service_slot.test.id
  target_resource_id  = azurerm_cosmosdb_account.test.id
  authentication {
    type         = "servicePrincipalSecret"
    client_id    = "someclientid"
    principal_id = azurerm_user_assigned_identity.test.principal_id
    secret       = "bar"
  }
}
`, template, data.RandomString, data.RandomInteger)
}

func (r AppServiceSlotConnectorResource) cosmosdbWithUserAssignedIdentity(data acceptance.TestData) string {
	template := r.templateWithCosmosDB(data)
	return fmt.Sprintf(`
%[1]s

data "azurerm_subscription" "test" {}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%[3]d"
  app_service_name    = azurerm_app_service.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  lifecycle {
    ignore_changes = [
      identity,
    ]
  }
}

resource "azurerm_app_service_slot_connection" "test" {
  name                = "acctestserviceconnector%[3]d"
  app_service_slot_id = azurerm_app_service_slot.test.id
  target_resource_id  = azurerm_cosmosdb_account.test.id
  authentication {
    type            = "userAssignedIdentity"
    subscription_id = data.azurerm_subscription.test.subscription_id
    client_id       = azurerm_user_assigned_identity.test.client_id
  }
}
`, template, data.RandomString, data.RandomInteger)
}

func (r AppServiceSlotConnectorResource) secretStore(data acceptance.TestData) string {
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
  name     = "acctestRG-sc-%[2]d"
  location = "%[1]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_key_vault" "test" {
  name                     = "accAKV-%[3]s"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  tenant_id                = data.azurerm_client_config.current.tenant_id
  sku_name                 = "standard"
  purge_protection_enabled = true
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
}

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%[2]d"
  app_service_name    = azurerm_app_service.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  lifecycle {
    ignore_changes = [
      identity,
    ]
  }
}

resource "azurerm_app_service_slot_connection" "test" {
  name                = "acctestserviceconnector%[2]d"
  app_service_slot_id = azurerm_app_service_slot.test.id
  target_resource_id  = azurerm_storage_account.test.id

  secret_store {
    key_vault_id = azurerm_key_vault.test.id
  }
  authentication {
    type = "systemAssignedIdentity"
  }
}
`, data.Locations.Primary, data.RandomInteger, data.RandomString)
}

func (r AppServiceSlotConnectorResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-sc-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceSlotConnectorResource) templateLinux(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-sc-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  os_type             = "Linux"
  sku_name            = "S1"
}

resource "azurerm_linux_web_app" "test" {
  name                = "acctestLWA-%d"
  location            = azurerm_resource_group.test.location
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
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceSlotConnectorResource) templateWithCosmosDB(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-sc-%[1]d"
  location = "%[2]s"
}

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctestacc%[3]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "GlobalDocumentDB"

  consistency_policy {
    consistency_level = "Session"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
