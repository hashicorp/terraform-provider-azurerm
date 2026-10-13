// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package containerapps_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/managedenvironmentsstorages"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ContainerAppEnvironmentStorageResource struct{}

func TestAccContainerAppEnvironmentStorage_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_environment_storage", "test")
	r := ContainerAppEnvironmentStorageResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("access_key"),
	})
}

func TestAccContainerAppEnvironmentStorage_nfsBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_environment_storage", "test")
	r := ContainerAppEnvironmentStorageResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.nfsBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerAppEnvironmentStorage_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_environment_storage", "test")
	r := ContainerAppEnvironmentStorageResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("access_key"),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("access_key"),
	})
}

func TestAccContainerAppEnvironmentStorage_keyVaultSecret(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_environment_storage", "test")
	r := ContainerAppEnvironmentStorageResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.keyVaultSecret(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerAppEnvironmentStorage_keyVaultSecretUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_environment_storage", "test")
	r := ContainerAppEnvironmentStorageResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.keyVaultAccessKey(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("access_key"),
		{
			Config: r.keyVaultSecret(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.keyVaultAccessKey(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("access_key"),
	})
}

func TestAccContainerAppEnvironmentStorage_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_environment_storage", "test")
	r := ContainerAppEnvironmentStorageResource{}

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

func (r ContainerAppEnvironmentStorageResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := managedenvironmentsstorages.ParseStorageID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.ContainerApps.StorageClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	if response.WasNotFound(resp.HttpResponse) {
		return pointer.To(false), nil
	}

	return pointer.To(true), nil
}

func (r ContainerAppEnvironmentStorageResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_container_app_environment_storage" "test" {
  name                         = "testacc-caes-%[2]d"
  container_app_environment_id = azurerm_container_app_environment.test.id
  account_name                 = azurerm_storage_account.test.name
  access_key                   = azurerm_storage_account.test.primary_access_key
  share_name                   = azurerm_storage_share.test.name
  access_mode                  = "ReadWrite"
}
`, r.template(data), data.RandomInteger)
}

func (r ContainerAppEnvironmentStorageResource) nfsBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_container_app_environment_storage" "test" {
  name                         = "testacc-caes-%[2]d"
  container_app_environment_id = azurerm_container_app_environment.test.id
  share_name                   = "/${azurerm_storage_account.test.name}/${azurerm_storage_share.test.name}"
  access_mode                  = "ReadWrite"
  nfs_server_url               = "${azurerm_storage_account.test.name}.file.core.windows.net"
}
`, r.template(data), data.RandomInteger)
}

func (r ContainerAppEnvironmentStorageResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_container_app_environment_storage" "test" {
  name                         = "testacc-caes-%[2]d"
  container_app_environment_id = azurerm_container_app_environment.test.id
  account_name                 = azurerm_storage_account.test.name
  access_key                   = azurerm_storage_account.test.secondary_access_key
  share_name                   = azurerm_storage_share.test.name
  access_mode                  = "ReadWrite"
}
`, r.template(data), data.RandomInteger)
}

func (r ContainerAppEnvironmentStorageResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_container_app_environment_storage" "import" {
  name                         = azurerm_container_app_environment_storage.test.name
  container_app_environment_id = azurerm_container_app_environment_storage.test.container_app_environment_id
  account_name                 = azurerm_container_app_environment_storage.test.account_name
  access_key                   = azurerm_container_app_environment_storage.test.access_key
  share_name                   = azurerm_container_app_environment_storage.test.share_name
  access_mode                  = azurerm_container_app_environment_storage.test.access_mode
}
`, r.basic(data))
}

func (r ContainerAppEnvironmentStorageResource) keyVaultSecret(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_container_app_environment_storage" "test" {
  name                         = "testacc-caes-%[2]d"
  container_app_environment_id = azurerm_container_app_environment.test.id
  account_name                 = azurerm_storage_account.test.name
  share_name                   = azurerm_storage_share.test.name
  access_mode                  = "ReadWrite"

  access_key_vault_secret {
    identity            = azurerm_user_assigned_identity.test.id
    key_vault_secret_id = azurerm_key_vault_secret.test.versionless_id
  }
}
`, r.keyVaultTemplate(data), data.RandomInteger)
}

func (r ContainerAppEnvironmentStorageResource) keyVaultAccessKey(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_container_app_environment_storage" "test" {
  name                         = "testacc-caes-%[2]d"
  container_app_environment_id = azurerm_container_app_environment.test.id
  account_name                 = azurerm_storage_account.test.name
  access_key                   = azurerm_storage_account.test.primary_access_key
  share_name                   = azurerm_storage_share.test.name
  access_mode                  = "ReadWrite"
}
`, r.keyVaultTemplate(data), data.RandomInteger)
}

func (r ContainerAppEnvironmentStorageResource) keyVaultTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy    = true
      recover_soft_deleted_key_vaults = true
    }
  }
}

%[1]s

data "azurerm_client_config" "current" {}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest-uai-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_key_vault" "test" {
  name                       = "acctest-kv-%[3]s"
  resource_group_name        = azurerm_resource_group.test.name
  location                   = azurerm_resource_group.test.location
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  rbac_authorization_enabled = false
  soft_delete_retention_days = 7

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id
    secret_permissions = [
      "Set",
      "Get",
      "Delete",
      "Purge",
      "Recover",
    ]
  }

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = azurerm_user_assigned_identity.test.principal_id
    secret_permissions = [
      "Get",
    ]
  }
}

resource "azurerm_key_vault_secret" "test" {
  name         = "storage-key-%[3]s"
  value        = azurerm_storage_account.test.primary_access_key
  key_vault_id = azurerm_key_vault.test.id
}

resource "azurerm_container_app_environment" "test" {
  name                       = "acctest-CAEnv%[2]d"
  resource_group_name        = azurerm_resource_group.test.name
  location                   = azurerm_resource_group.test.location
  log_analytics_workspace_id = azurerm_log_analytics_workspace.test.id
  logs_destination           = "log-analytics"

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
}
`, r.templateBase(data), data.RandomInteger, data.RandomString)
}

func (r ContainerAppEnvironmentStorageResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[2]s

resource "azurerm_container_app_environment" "test" {
  name                       = "acctest-CAEnv%[1]d"
  resource_group_name        = azurerm_resource_group.test.name
  location                   = azurerm_resource_group.test.location
  log_analytics_workspace_id = azurerm_log_analytics_workspace.test.id
  logs_destination           = "log-analytics"
}
`, data.RandomInteger, r.templateBase(data))
}

func (r ContainerAppEnvironmentStorageResource) templateBase(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-CAE-%[1]d"
  location = "%[2]s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%[3]s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "accTest"
  }
}

resource "azurerm_storage_share" "test" {
  name               = "testshare%[3]s"
  storage_account_id = azurerm_storage_account.test.id
  quota              = 1
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
