---
subcategory: "Synapse"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_synapse_workspace"
description: |-
  Manages a Synapse Workspace.
---

# azurerm_synapse_workspace

Manages a Synapse Workspace.

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
```

## Example Usage - creating a workspace with Customer Managed Key and Azure AD Admin

```hcl
data "azurerm_client_config" "current" {}

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
    "Create", "Get", "Delete", "Purge", "GetRotationPolicy"
  ]
}

resource "azurerm_key_vault_key" "example" {
  name         = "workspaceencryptionkey"
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

  customer_managed_key {
    key_versionless_id = azurerm_key_vault_key.example.versionless_id
    key_name           = "enckey"
  }

  identity {
    type = "SystemAssigned"
  }

  tags = {
    Env = "production"
  }
}

resource "azurerm_key_vault_access_policy" "workspace_policy" {
  key_vault_id = azurerm_key_vault.example.id
  tenant_id    = azurerm_synapse_workspace.example.identity[0].tenant_id
  object_id    = azurerm_synapse_workspace.example.identity[0].principal_id

  key_permissions = [
    "Get", "WrapKey", "UnwrapKey"
  ]
}

resource "azurerm_synapse_workspace_key" "example" {
  customer_managed_key_versionless_id = azurerm_key_vault_key.example.versionless_id
  synapse_workspace_id                = azurerm_synapse_workspace.example.id
  active                              = true
  customer_managed_key_name           = "enckey"
  depends_on                          = [azurerm_key_vault_access_policy.workspace_policy]
}

resource "azurerm_synapse_workspace_aad_admin" "example" {
  synapse_workspace_id = azurerm_synapse_workspace.example.id
  login                = "AzureAD Admin"
  object_id            = "00000000-0000-0000-0000-000000000000"
  tenant_id            = "00000000-0000-0000-0000-000000000000"

  depends_on = [azurerm_synapse_workspace_key.example]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this synapse Workspace. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the synapse Workspace should exist. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the Azure Region where the synapse Workspace should exist. Changing this forces a new resource to be created.

* `identity` - (Optional) An `identity` block as defined below.

* `storage_data_lake_gen2_filesystem_id` - (Required) Specifies the ID of storage data lake gen2 filesystem resource. Changing this forces a new resource to be created.

* `sql_administrator_login` - (Optional) Specifies The login name of the SQL administrator. Changing this forces a new resource to be created. If this is not provided `customer_managed_key` must be provided.

* `sql_administrator_login_password` - (Optional) The Password associated with the `sql_administrator_login` for the SQL administrator. If this is not provided `customer_managed_key` must be provided.

* `azuread_authentication_only` - (Optional) Is Azure Active Directory Authentication the only way to authenticate with resources inside this synapse Workspace. Defaults to `false`.

---

* `compute_subnet_id` - (Optional) Subnet ID used for computes in workspace Changing this forces a new resource to be created.

* `azure_devops_repo` - (Optional) An `azure_devops_repo` block as defined below.

* `data_exfiltration_protection_enabled` - (Optional) Is data exfiltration protection enabled in this workspace? If set to `true`, `managed_virtual_network_enabled` must also be set to `true`. Changing this forces a new resource to be created.

* `customer_managed_key` - (Optional) A `customer_managed_key` block as defined below.

* `github_repo` - (Optional) A `github_repo` block as defined below.

* `linking_allowed_for_aad_tenant_ids` - (Optional) Allowed AAD Tenant Ids For Linking.

* `managed_resource_group_name` - (Optional) Workspace managed resource group. Changing this forces a new resource to be created.

* `managed_virtual_network_enabled` - (Optional) Is Virtual Network enabled for all computes in this workspace? Changing this forces a new resource to be created.

* `public_network_access_enabled` - (Optional) Whether public network access is allowed for the Cognitive Account. Defaults to `true`.

* `purview_id` - (Optional) The ID of purview account.

* `sql_identity_control_enabled` - (Optional) Are pipelines (running as workspace's system assigned identity) allowed to access SQL pools?

* `tags` - (Optional) A mapping of tags which should be assigned to the Synapse Workspace.

---

An `azure_devops_repo` block supports the following:

* `account_name` - (Required) Specifies the Azure DevOps account name.

* `branch_name` - (Required) Specifies the collaboration branch of the repository to get code from.

* `last_commit_id` - (Optional) The last commit ID.

* `project_name` - (Required) Specifies the name of the Azure DevOps project.

* `repository_name` - (Required) Specifies the name of the git repository.

* `root_folder` - (Required) Specifies the root folder within the repository. Set to `/` for the top level.

* `tenant_id` - (Optional) the ID of the tenant for the Azure DevOps account.

---

A `customer_managed_key` block supports the following:

* `key_versionless_id` - (Required) The Azure Key Vault Key Versionless ID to be used as the Customer Managed Key (CMK) for double encryption (e.g. `https://example-keyvault.vault.azure.net/type/cmk/`).

* `key_name` - (Optional) An identifier for the key. Name needs to match the name of the key used with the `azurerm_synapse_workspace_key` resource. Defaults to "cmk" if not specified.

* `user_assigned_identity_id` - (Optional) The User Assigned Identity ID to be used for accessing the Customer Managed Key for encryption.

---

The `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be associated with this Synapse Workspace. Possible values are `SystemAssigned`, `UserAssigned` and `SystemAssigned, UserAssigned` (to enable both).

* `identity_ids` - (Optional) Specifies a list of User Assigned Managed Identity IDs to be assigned to this Synapse Workspace.

~> **Note:** This is required when `type` is set to `UserAssigned` or `SystemAssigned, UserAssigned`.

---

A `github_repo` block supports the following:

* `account_name` - (Required) Specifies the GitHub account name.

* `branch_name` - (Required) Specifies the collaboration branch of the repository to get code from.

* `last_commit_id` - (Optional) The last commit ID.

* `repository_name` - (Required) Specifies the name of the git repository.

* `root_folder` - (Required) Specifies the root folder within the repository. Set to `/` for the top level.

* `git_url` - (Optional) Specifies the GitHub Enterprise host name. For example: <https://github.mydomain.com>.

-> **Note:** You must log in to the Synapse UI to complete the authentication to the GitHub repository.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the synapse Workspace.

* `connectivity_endpoints` - A map of Connectivity endpoints for this Synapse Workspace. Possible key values are `dev`, `sql`, `sqlOnDemand`, and `web`.

---

The `identity` block exports the following:

* `principal_id` - The Principal ID for the Service Principal associated with the Managed Service Identity of this Synapse Workspace.

* `tenant_id` - The Tenant ID for the Service Principal associated with the Managed Service Identity of this Synapse Workspace.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Synapse Workspace.
* `read` - (Defaults to 5 minutes) Used when retrieving the Synapse Workspace.
* `update` - (Defaults to 30 minutes) Used when updating the Synapse Workspace.
* `delete` - (Defaults to 30 minutes) Used when deleting the Synapse Workspace.

## Import

Synapse Workspace can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_synapse_workspace.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Synapse/workspaces/workspace1
```
