// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package synapse_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SynapseWorkspaceKeysResource struct{}

func TestAccSynapseWorkspaceKeys_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_workspace_key", "test")
	r := SynapseWorkspaceKeysResource{}

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

func TestAccSynapseWorkspaceKeys_basicUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_workspace_key", "test")
	r := SynapseWorkspaceKeysResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r SynapseWorkspaceKeysResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.WorkspaceKeysID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Synapse.KeysClient.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.KeyName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Synapse Workspace Key %q (Workspace %q): %+v", id.KeyName, id.WorkspaceName, err)
	}

	return utils.Bool(true), nil
}

func (r SynapseWorkspaceKeysResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_synapse_workspace" "test" {
  name                                 = "acctestsw%[2]d"
  resource_group_name                  = azurerm_resource_group.test.name
  location                             = azurerm_resource_group.test.location
  storage_data_lake_gen2_filesystem_id = azurerm_storage_data_lake_gen2_filesystem.test.id
  sql_administrator_login              = "sqladminuser"
  sql_administrator_login_password     = "H@Sh1CoR3!"
  customer_managed_key {
    key_versionless_id = azurerm_key_vault_key.test.versionless_id
    key_name           = "test_key"
  }
  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_key_vault_access_policy" "workspace_policy" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = azurerm_synapse_workspace.test.identity[0].tenant_id
  object_id    = azurerm_synapse_workspace.test.identity[0].principal_id

  key_permissions = [
    "Get", "WrapKey", "UnwrapKey"
  ]
}

resource "azurerm_synapse_workspace_key" "test" {
  customer_managed_key_versionless_id = azurerm_key_vault_key.test.versionless_id
  synapse_workspace_id                = azurerm_synapse_workspace.test.id
  active                              = true
  customer_managed_key_name           = "test_key"
  depends_on                          = [azurerm_key_vault_access_policy.workspace_policy]
}

resource "azurerm_synapse_workspace_key" "test2" {
  customer_managed_key_versionless_id = azurerm_key_vault_key.test2.versionless_id
  synapse_workspace_id                = azurerm_synapse_workspace.test.id
  active                              = false
  customer_managed_key_name           = "test_key2"
  depends_on                          = [azurerm_key_vault_access_policy.workspace_policy, azurerm_synapse_workspace_key.test]
}
`, template, data.RandomInteger)
}

func (r SynapseWorkspaceKeysResource) basicUpdate(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_synapse_workspace" "test" {
  name                                 = "acctestsw%[2]d"
  resource_group_name                  = azurerm_resource_group.test.name
  location                             = azurerm_resource_group.test.location
  storage_data_lake_gen2_filesystem_id = azurerm_storage_data_lake_gen2_filesystem.test.id
  sql_administrator_login              = "sqladminuser"
  sql_administrator_login_password     = "H@Sh1CoR3!"
  customer_managed_key {
    key_versionless_id = azurerm_key_vault_key.test2.versionless_id
    key_name           = "test_key2"
  }
  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_key_vault_access_policy" "workspace_policy" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = azurerm_synapse_workspace.test.identity[0].tenant_id
  object_id    = azurerm_synapse_workspace.test.identity[0].principal_id

  key_permissions = [
    "Get", "WrapKey", "UnwrapKey"
  ]
}

resource "azurerm_synapse_workspace_key" "test" {
  customer_managed_key_versionless_id = azurerm_key_vault_key.test.versionless_id
  synapse_workspace_id                = azurerm_synapse_workspace.test.id
  active                              = false
  customer_managed_key_name           = "test_key"
  depends_on                          = [azurerm_key_vault_access_policy.workspace_policy]
}

resource "azurerm_synapse_workspace_key" "test2" {
  customer_managed_key_versionless_id = azurerm_key_vault_key.test2.versionless_id
  synapse_workspace_id                = azurerm_synapse_workspace.test.id
  active                              = true
  customer_managed_key_name           = "test_key2"
  depends_on                          = [azurerm_key_vault_access_policy.workspace_policy, azurerm_synapse_workspace_key.test]
}
`, template, data.RandomInteger)
}

func (r SynapseWorkspaceKeysResource) template(data acceptance.TestData) string {
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
  name     = "acctestRG-synapse-%[3]d"
  location = "%[1]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestacc%[2]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_kind             = "BlobStorage"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_data_lake_gen2_filesystem" "test" {
  name               = "acctest-%[3]d"
  storage_account_id = azurerm_storage_account.test.id
}

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "test" {
  name                     = "acckv%[3]d"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  tenant_id                = data.azurerm_client_config.current.tenant_id
  sku_name                 = "standard"
  purge_protection_enabled = true
}

resource "azurerm_key_vault_access_policy" "deployer" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Create", "Get", "Delete", "Purge", "GetRotationPolicy"
  ]
}

resource "azurerm_key_vault_key" "test" {
  name         = "key"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts = [
    "unwrapKey",
    "wrapKey"
  ]
  depends_on = [
    azurerm_key_vault_access_policy.deployer
  ]
}

resource "azurerm_key_vault_key" "test2" {
  name         = "key2"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts = [
    "unwrapKey",
    "wrapKey"
  ]
  depends_on = [
    azurerm_key_vault_access_policy.deployer
  ]
}

`, data.Locations.Primary, data.RandomString, data.RandomInteger)
}
