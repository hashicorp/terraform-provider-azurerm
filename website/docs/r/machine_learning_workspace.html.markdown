---
subcategory: "Machine Learning"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_machine_learning_workspace"
description: |-
  Manages a Azure Machine Learning Workspace.
---
# azurerm_machine_learning_workspace

Manages a Azure Machine Learning Workspace

~> **Note:** For examples on how to set up the Azure Machine Learning workspace, together with compute and integrated services, see [Terraform Quickstart](https://github.com/Azure/terraform/tree/master/quickstart)

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_application_insights" "example" {
  name                = "workspace-example-ai"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  application_type    = "web"
}

resource "azurerm_key_vault" "example" {
  name                = "workspaceexamplekeyvault"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "premium"
}

resource "azurerm_storage_account" "example" {
  name                     = "workspacestorageaccount"
  location                 = azurerm_resource_group.example.location
  resource_group_name      = azurerm_resource_group.example.name
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_machine_learning_workspace" "example" {
  name                    = "example-workspace"
  location                = azurerm_resource_group.example.location
  resource_group_name     = azurerm_resource_group.example.name
  application_insights_id = azurerm_application_insights.example.id
  key_vault_id            = azurerm_key_vault.example.id
  storage_account_id      = azurerm_storage_account.example.id

  identity {
    type = "SystemAssigned"
  }
}
```

## Example Usage with Data encryption

~> **Note:** The Key Vault must enable purge protection.

```hcl
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy = false
    }
  }
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_application_insights" "example" {
  name                = "workspace-example-ai"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  application_type    = "web"
}

resource "azurerm_key_vault" "example" {
  name                     = "workspaceexamplekeyvault"
  location                 = azurerm_resource_group.example.location
  resource_group_name      = azurerm_resource_group.example.name
  tenant_id                = data.azurerm_client_config.current.tenant_id
  sku_name                 = "premium"
  purge_protection_enabled = true
}
resource "azurerm_key_vault_access_policy" "example" {
  key_vault_id = azurerm_key_vault.example.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Create",
    "Get",
    "Delete",
    "Purge",
    "GetRotationPolicy",
  ]
}

resource "azurerm_storage_account" "example" {
  name                     = "workspacestorageaccount"
  location                 = azurerm_resource_group.example.location
  resource_group_name      = azurerm_resource_group.example.name
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_key_vault_key" "example" {
  name         = "workspaceexamplekeyvaultkey"
  key_vault_id = azurerm_key_vault.example.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]

  depends_on = [azurerm_key_vault.example, azurerm_key_vault_access_policy.example]
}

resource "azurerm_machine_learning_workspace" "example" {
  name                    = "example-workspace"
  location                = azurerm_resource_group.example.location
  resource_group_name     = azurerm_resource_group.example.name
  application_insights_id = azurerm_application_insights.example.id
  key_vault_id            = azurerm_key_vault.example.id
  storage_account_id      = azurerm_storage_account.example.id

  identity {
    type = "SystemAssigned"
  }

  encryption {
    key_vault_id = azurerm_key_vault.example.id
    key_id       = azurerm_key_vault_key.example.id
  }
}
```

## Example Usage with User Assigned Identity and Data Encryption

~> **Note:** The Key Vault must enable purge protection.

```hcl
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy = false
    }
  }
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_application_insights" "example" {
  name                = "example-ai"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  application_type    = "web"
}

resource "azurerm_storage_account" "example" {
  name                     = "examplestorageaccount"
  location                 = azurerm_resource_group.example.location
  resource_group_name      = azurerm_resource_group.example.name
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_key_vault" "example" {
  name                     = "example-keyvalut"
  location                 = azurerm_resource_group.example.location
  resource_group_name      = azurerm_resource_group.example.name
  tenant_id                = data.azurerm_client_config.current.tenant_id
  sku_name                 = "premium"
  purge_protection_enabled = true
}

resource "azurerm_user_assigned_identity" "example" {
  name                = "example-identity"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_key_vault_access_policy" "example-identity" {
  key_vault_id = azurerm_key_vault.example.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = azurerm_user_assigned_identity.example.principal_id

  // default set by service
  key_permissions = [
    "WrapKey",
    "UnwrapKey",
    "Get",
    "Recover",
  ]


  secret_permissions = [
    "Get",
    "List",
    "Set",
    "Delete",
    "Recover",
    "Backup",
    "Restore"
  ]
}

resource "azurerm_key_vault_access_policy" "example-sp" {
  key_vault_id = azurerm_key_vault.example.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Get",
    "Create",
    "Recover",
    "Delete",
    "Purge",
    "GetRotationPolicy",
  ]
}

data "azuread_service_principal" "test" {
  display_name = "Azure Cosmos DB"
}

resource "azurerm_key_vault_access_policy" "example-cosmosdb" {
  key_vault_id = azurerm_key_vault.example.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azuread_service_principal.test.object_id

  key_permissions = [
    "Get",
    "Recover",
    "UnwrapKey",
    "WrapKey",
  ]
  depends_on = [data.azuread_service_principal.test, data.azurerm_client_config.current]
}

resource "azurerm_key_vault_key" "example" {
  name         = "example-keyvaultkey"
  key_vault_id = azurerm_key_vault.example.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]
  depends_on = [azurerm_key_vault.example, azurerm_key_vault_access_policy.example-sp]
}

resource "azurerm_role_assignment" "example-role1" {
  scope                = azurerm_key_vault.example.id
  role_definition_name = "Contributor"
  principal_id         = azurerm_user_assigned_identity.example.principal_id
}

resource "azurerm_role_assignment" "example-role2" {
  scope                = azurerm_storage_account.example.id
  role_definition_name = "Storage Blob Data Contributor"
  principal_id         = azurerm_user_assigned_identity.example.principal_id
}

resource "azurerm_role_assignment" "example-role3" {
  scope                = azurerm_storage_account.example.id
  role_definition_name = "Contributor"
  principal_id         = azurerm_user_assigned_identity.example.principal_id
}

resource "azurerm_role_assignment" "example-role4" {
  scope                = azurerm_application_insights.example.id
  role_definition_name = "Contributor"
  principal_id         = azurerm_user_assigned_identity.example.principal_id
}


resource "azurerm_machine_learning_workspace" "example" {
  name                    = "example-workspace"
  location                = azurerm_resource_group.example.location
  resource_group_name     = azurerm_resource_group.example.name
  application_insights_id = azurerm_application_insights.example.id
  key_vault_id            = azurerm_key_vault.example.id
  storage_account_id      = azurerm_storage_account.example.id

  high_business_impact = true

  primary_user_assigned_identity = azurerm_user_assigned_identity.example.id
  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.example.id,
    ]
  }

  encryption {
    user_assigned_identity_id = azurerm_user_assigned_identity.example.id
    key_vault_id              = azurerm_key_vault.example.id
    key_id                    = azurerm_key_vault_key.example.id
  }

  depends_on = [
    azurerm_role_assignment.example-role1, azurerm_role_assignment.example-role2, azurerm_role_assignment.example-role3,
    azurerm_role_assignment.example-role4,
    azurerm_key_vault_access_policy.example-cosmosdb,
  ]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Machine Learning Workspace. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group in which the Machine Learning Workspace should exist. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the Machine Learning Workspace should exist. Changing this forces a new resource to be created.

* `application_insights_id` - (Required) The ID of the Application Insights associated with this Machine Learning Workspace. Changing this forces a new resource to be created.

* `key_vault_id` - (Required) The ID of key vault associated with this Machine Learning Workspace. Changing this forces a new resource to be created.

* `storage_account_id` - (Required) The ID of the Storage Account associated with this Machine Learning Workspace. Changing this forces a new resource to be created.

-> **Note:** The `account_tier` cannot be `Premium` in order to associate the Storage Account to this Machine Learning Workspace.

* `identity` - (Required) An `identity` block as defined below.

* `kind` - (Optional) The type of the Workspace. Possible values are `Default`, `FeatureStore`. Defaults to `Default`

* `container_registry_id` - (Optional) The ID of the container registry associated with this Machine Learning Workspace. Changing this forces a new resource to be created.

-> **Note:** The `admin_enabled` should be `true` in order to associate the Container Registry to this Machine Learning Workspace.

* `public_network_access_enabled` - (Optional) Enable public access when this Machine Learning Workspace is behind VNet. Defaults to `true`.

~> **Note:** `public_access_behind_virtual_network_enabled` is deprecated and will be removed in favour of the property `public_network_access_enabled`.

* `image_build_compute_name` - (Optional) The compute name for image build of the Machine Learning Workspace.

* `description` - (Optional) The description of this Machine Learning Workspace.

* `encryption` - (Optional) An `encryption` block as defined below. Changing this forces a new resource to be created.

* `managed_network` - (Optional) A `managed_network` block as defined below.

* `feature_store` - (Optional) A `feature_store` block as defined below.

* `friendly_name` - (Optional) Display name for this Machine Learning Workspace.

* `high_business_impact` - (Optional) Flag to signal High Business Impact (HBI) data in the workspace and reduce diagnostic data collected by the service. Changing this forces a new resource to be created.

* `primary_user_assigned_identity` - (Optional) The user assigned identity id that represents the workspace identity.

* `v1_legacy_mode_enabled` - (Optional) Enable V1 API features, enabling `v1_legacy_mode` may prevent you from using features provided by the v2 API. Defaults to `false`.

* `sku_name` - (Optional) SKU/edition of the Machine Learning Workspace, possible values are `Free`, `Basic`, `Standard` and `Premium`. Defaults to `Basic`.

* `serverless_compute` - (Optional) A `serverless_compute` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Machine Learning Workspace. Possible values are `SystemAssigned`, `UserAssigned`, `SystemAssigned, UserAssigned` (to enable both).

* `identity_ids` - (Optional) Specifies a list of User Assigned Managed Identity IDs to be assigned to this Machine Learning Workspace.

~> **Note:** This is required when `type` is set to `UserAssigned` or `SystemAssigned, UserAssigned`.

---

An `encryption` block supports the following:

* `key_vault_id` - (Required) The ID of the keyVault where the customer owned encryption key is present.

* `key_id` - (Required) The Key Vault URI to access the encryption key.

* `user_assigned_identity_id` - (Optional) The Key Vault URI to access the encryption key.

~> **Note:** `user_assigned_identity_id` must set when`identity.type` is `UserAssigned` or service won't be able to find the assigned permissions.

---

An `managed_network` block supports the following:

* `isolation_mode` - (Optional) The isolation mode of the Machine Learning Workspace. Possible values are `Disabled`, `AllowOnlyApprovedOutbound`, and `AllowInternetOutbound`

---

A `serverless_compute` block supports the following:

* `subnet_id` - (Optional) The ID of an existing Virtual Network Subnet in which the serverless compute nodes should be deployed to.

* `public_ip_enabled` - (Optional) Should serverless compute nodes deployed in a custom Virtual Network have public IP addresses enabled for a workspace with private endpoint? Defaults to `false`.

~> **Note:** `public_ip_enabled` cannot be updated from `true` to `false` when `subnet_id` is not set. `public_ip_enabled` must be set to `true` if `subnet_id` is not set and when `public_network_access_enabled` is `false`.

---

An `feature_store` block supports the following:

* `computer_spark_runtime_version` - (Optional) The version of Spark runtime.

* `offline_connection_name` - (Optional) The name of offline store connection.

* `online_connection_name` - (Optional) The name of online store connection.

~> **Note:** `feature_store` must be set when`kind` is `FeatureStore`


## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Machine Learning Workspace.

* `discovery_url` - The url for the discovery service to identify regional endpoints for machine learning experimentation services.

* `workspace_id` - The immutable id associated with this workspace.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Machine Learning Workspace.
* `read` - (Defaults to 5 minutes) Used when retrieving the Machine Learning Workspace.
* `update` - (Defaults to 30 minutes) Used when updating the Machine Learning Workspace.
* `delete` - (Defaults to 30 minutes) Used when deleting the Machine Learning Workspace.

## Import

Machine Learning Workspace can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_machine_learning_workspace.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.MachineLearningServices/workspaces/workspace1
```
