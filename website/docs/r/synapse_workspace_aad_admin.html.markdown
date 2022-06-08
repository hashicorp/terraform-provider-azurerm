---
subcategory: "Synapse"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_synapse_workspace_aad_admin"
description: |-
  Manages Synapse Workspace AAD Admin
---

# azurerm_synapse_workspace_aad_admin

Manages an Azure Active Directory Administrator setting for a Synapse Workspace

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "examplestorageacc"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
  account_kind             = "StorageV2"
  is_hns_enabled           = "true"
}

resource "azurerm_storage_data_lake_gen2_filesystem" "example" {
  name               = "example"
  storage_account_id = azurerm_storage_account.example.id
}

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "example" {
  name                     = "example"
  location                 = azurerm_resource_group.example.location
  resource_group_name      = azurerm_resource_group.example.name
  tenant_id                = data.azurerm_client_config.current.tenant_id
  sku_name                 = "standard"
  purge_protection_enabled = true
}

resource "azurerm_key_vault_access_policy" "deployer" {
  key_vault_id = azurerm_key_vault.example.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Create", "Get", "Delete", "Purge"
  ]
}

resource "azurerm_key_vault_key" "example" {
  name         = "workspace-encryption-key"
  key_vault_id = azurerm_key_vault.example.id
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

resource "azurerm_synapse_workspace" "example" {
  name                                 = "example"
  resource_group_name                  = azurerm_resource_group.example.name
  location                             = azurerm_resource_group.example.location
  storage_data_lake_gen2_filesystem_id = azurerm_storage_data_lake_gen2_filesystem.example.id
  sql_administrator_login              = "sqladminuser"
  sql_administrator_login_password     = "H@Sh1CoR3!"

  identity {
    type = "SystemAssigned"
  }

  tags = {
    Env = "production"
  }
}

resource "azurerm_synapse_workspace_aad_admin" "example" {
  synapse_workspace_id = azurerm_synapse_workspace.example.id
  login                = "AzureAD Admin"
  object_id            = data.azurerm_client_config.current.object_id
  tenant_id            = data.azurerm_client_config.current.tenant_id
}
```

## Arguments Reference

The following arguments are supported:

* `synapse_workspace_id` - (Required) The ID of the Synapse Workspace where the Azure AD Administrator should be configured. 

* `login` - (Required) The login name of the Azure AD Administrator of this Synapse Workspace.

* `object_id` - (Required) The object id of the Azure AD Administrator of this Synapse Workspace.

* `tenant_id` - (Required) The tenant id of the Azure AD Administrator of this Synapse Workspace.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Synapse Workspace.
* `read` - (Defaults to 5 minutes) Used when retrieving the Synapse Workspace.
* `update` - (Defaults to 30 minutes) Used when updating the Synapse Workspace.
* `delete` - (Defaults to 30 minutes) Used when deleting the Synapse Workspace.

## Import

Synapse Workspace Azure AD Administrator can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_synapse_workspace_aad_admin.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Synapse/workspaces/workspace1/administrators/activeDirectory
```
