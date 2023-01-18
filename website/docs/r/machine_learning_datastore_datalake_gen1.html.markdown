---
subcategory: "Machine Learning"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_machine_learning_datastore_datalake_gen1"
description: |-
  Manages a Machine Learning Data Lake Gen1 DataStore.
---

# azurerm_machine_learning_datastore_datalake_gen1

Manages a Machine Learning Data Lake Gen1 DataStore.

## Example Usage with Data Lake Gen1

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

resource "azurerm_machine_learning_datastore_datalake_gen1" "example" {
  name         = "example-datastore"
  workspace_id = azurerm_machine_learning_workspace.example.id
  store_name   = "azure-datalake-gen1"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Machine Learning DataStore. Changing this forces a new Machine Learning DataStore to be created.

* `workspace_id` - (Required) The ID of the Machine Learning Workspace. Changing this forces a new Machine Learning DataStore to be created.

* `store_name` - (Required) The name of the Azure Data Lake Gen1. Changing this forces a new Machine Learning DataStore to be created.

---
* `tenant_id` - (Optional) The ID of the Tenant which the Service Principal belongs to.

* `client_id` - (Optional) The object ID of the Service Principal.

* `client_secret` - (Optional) The secret of the Service Principal.

* `authority_url` - (Optional) An URL used for authentication.

* `description` - (Optional) Text used to describe the asset. Changing this forces a new Machine Learning DataStore to be created.

* `is_default` - (Optional) Specifies whether this Machines Learning DataStore is the default for the Workspace. Defaults to `false`.

~> **Note:** `is_default` can only be set to `true` on update.

* `service_data_auth_identity` - (Optional) Specifies which identity to use when retrieving data from the specified source. Defaults to `None`. Possible values are `None`, `WorkspaceSystemAssignedIdentity` and `WorkspaceUserAssignedIdentity`.

* `tags` - (Optional) A mapping of tags which should be assigned to the Machine Learning DataStore. Changing this forces a new Machine Learning DataStore to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Machine Learning DataStore.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Machine Learning DataStore.
* `read` - (Defaults to 5 minutes) Used when retrieving the Machine Learning DataStore.
* `update` - (Defaults to 30 minutes) Used when updating the Machine Learning DataStore.
* `delete` - (Defaults to 30 minutes) Used when deleting the Machine Learning DataStore.

## Import

Machine Learning DataStores can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_machine_learning_datastore_datalake_gen1.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.MachineLearningServices/workspaces/mlw1/datastores/datastore1
```
