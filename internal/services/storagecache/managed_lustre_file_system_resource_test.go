// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storagecache_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2023-05-01/amlfilesystems"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ManagedLustreFileSystemResource struct{}

func TestAccManagedLustreFileSystem_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_lustre_file_system", "test")
	r := ManagedLustreFileSystemResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("mgs_address").IsNotEmpty(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedLustreFileSystem_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_lustre_file_system", "test")
	r := ManagedLustreFileSystemResource{}

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

func TestAccManagedLustreFileSystem_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_lustre_file_system", "test")
	r := ManagedLustreFileSystemResource{}

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

func TestAccManagedLustreFileSystem_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_lustre_file_system", "test")
	r := ManagedLustreFileSystemResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r ManagedLustreFileSystemResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := amlfilesystems.ParseAmlFilesystemID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.StorageCache.AmlFilesystems
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r ManagedLustreFileSystemResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-amlfs-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctest-subnet-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r ManagedLustreFileSystemResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_managed_lustre_file_system" "test" {
  name                   = "acctest-amlfs-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  sku_name               = "AMLFS-Durable-Premium-250"
  subnet_id              = azurerm_subnet.test.id
  storage_capacity_in_tb = 8
  zones                  = ["1"]

  maintenance_window {
    day_of_week        = "Friday"
    time_of_day_in_utc = "22:00"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ManagedLustreFileSystemResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_managed_lustre_file_system" "import" {
  name                   = azurerm_managed_lustre_file_system.test.name
  resource_group_name    = azurerm_managed_lustre_file_system.test.resource_group_name
  location               = azurerm_managed_lustre_file_system.test.location
  sku_name               = azurerm_managed_lustre_file_system.test.sku_name
  subnet_id              = azurerm_managed_lustre_file_system.test.subnet_id
  storage_capacity_in_tb = azurerm_managed_lustre_file_system.test.storage_capacity_in_tb
  zones                  = azurerm_managed_lustre_file_system.test.zones

  maintenance_window {
    day_of_week        = "Friday"
    time_of_day_in_utc = "22:00"
  }
}
`, r.basic(data))
}

func (r ManagedLustreFileSystemResource) templateForComplete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_client_config" "current" {}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestmi%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_key_vault" "test" {
  name                     = "acctestkv%s"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  tenant_id                = data.azurerm_client_config.current.tenant_id
  sku_name                 = "standard"
  purge_protection_enabled = true
}

resource "azurerm_key_vault_access_policy" "server" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = azurerm_user_assigned_identity.test.principal_id

  key_permissions = ["Get", "List", "WrapKey", "UnwrapKey", "GetRotationPolicy", "SetRotationPolicy"]
}

resource "azurerm_key_vault_access_policy" "client" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = ["Get", "Create", "Delete", "List", "Restore", "Recover", "UnwrapKey", "WrapKey", "Purge", "Encrypt", "Decrypt", "Sign", "Verify", "GetRotationPolicy", "SetRotationPolicy"]
}

resource "azurerm_key_vault_key" "test" {
  name         = "test"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]

  depends_on = [
    azurerm_key_vault_access_policy.client,
    azurerm_key_vault_access_policy.server,
  ]
}

resource "azurerm_storage_account" "test" {
  name                            = "acctestsa%s"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  account_tier                    = "Standard"
  account_replication_type        = "LRS"
  allow_nested_items_to_be_public = true
}

resource "azurerm_storage_container" "test" {
  name                  = "storagecontainer"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_storage_container" "test2" {
  name                  = "storagecontainer2"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

data "azuread_service_principal" "test" {
  display_name = "HPC Cache Resource Provider"
}

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_storage_account.test.id
  role_definition_name = "Contributor"
  principal_id         = data.azuread_service_principal.test.object_id
}

resource "azurerm_role_assignment" "test2" {
  scope                = azurerm_storage_account.test.id
  role_definition_name = "Storage Blob Data Contributor"
  principal_id         = data.azuread_service_principal.test.object_id
}
`, r.template(data), data.RandomInteger, data.RandomString, data.RandomString)
}

func (r ManagedLustreFileSystemResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_managed_lustre_file_system" "test" {
  name                   = "acctest-amlfs-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  sku_name               = "AMLFS-Durable-Premium-250"
  subnet_id              = azurerm_subnet.test.id
  storage_capacity_in_tb = 8
  zones                  = ["1"]

  maintenance_window {
    day_of_week        = "Friday"
    time_of_day_in_utc = "22:00"
  }

  identity {
    type = "UserAssigned"

    identity_ids = [
      azurerm_user_assigned_identity.test.id
    ]
  }

  encryption_key {
    key_url         = azurerm_key_vault_key.test.id
    source_vault_id = azurerm_key_vault.test.id
  }

  hsm_setting {
    container_id         = azurerm_storage_container.test.resource_manager_id
    logging_container_id = azurerm_storage_container.test2.resource_manager_id
    import_prefix        = "/"
  }

  tags = {
    Env = "Test"
  }

  depends_on = [azurerm_role_assignment.test, azurerm_role_assignment.test2]
}
`, r.templateForComplete(data), data.RandomInteger)
}

func (r ManagedLustreFileSystemResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault" "test2" {
  name                     = "acctestkv2%s"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  tenant_id                = data.azurerm_client_config.current.tenant_id
  sku_name                 = "standard"
  purge_protection_enabled = true
}

resource "azurerm_key_vault_access_policy" "server2" {
  key_vault_id = azurerm_key_vault.test2.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = azurerm_user_assigned_identity.test.principal_id

  key_permissions = ["Get", "List", "WrapKey", "UnwrapKey", "GetRotationPolicy", "SetRotationPolicy"]
}

resource "azurerm_key_vault_access_policy" "client2" {
  key_vault_id = azurerm_key_vault.test2.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = ["Get", "Create", "Delete", "List", "Restore", "Recover", "UnwrapKey", "WrapKey", "Purge", "Encrypt", "Decrypt", "Sign", "Verify", "GetRotationPolicy", "SetRotationPolicy"]
}

resource "azurerm_key_vault_key" "test2" {
  name         = "test2"
  key_vault_id = azurerm_key_vault.test2.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]

  depends_on = [
    azurerm_key_vault_access_policy.client2,
    azurerm_key_vault_access_policy.server2,
  ]
}

resource "azurerm_managed_lustre_file_system" "test" {
  name                   = "acctest-amlfs-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  sku_name               = "AMLFS-Durable-Premium-250"
  subnet_id              = azurerm_subnet.test.id
  storage_capacity_in_tb = 8
  zones                  = ["1"]

  maintenance_window {
    day_of_week        = "Thursday"
    time_of_day_in_utc = "21:00"
  }

  identity {
    type = "UserAssigned"

    identity_ids = [
      azurerm_user_assigned_identity.test.id
    ]
  }

  encryption_key {
    key_url         = azurerm_key_vault_key.test2.id
    source_vault_id = azurerm_key_vault.test2.id
  }

  hsm_setting {
    container_id         = azurerm_storage_container.test.resource_manager_id
    logging_container_id = azurerm_storage_container.test2.resource_manager_id
    import_prefix        = "/"
  }

  tags = {
    Env = "Test2"
  }

  depends_on = [azurerm_role_assignment.test, azurerm_role_assignment.test2]
}
`, r.templateForComplete(data), data.RandomString, data.RandomInteger)
}
