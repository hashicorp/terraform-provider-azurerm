---
subcategory: "Machine Learning"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_machine_learning_datastore"
description: |-
  Manages a Machine Learning DataStore.
---

# azurerm_machine_learning_datastore

Manages a Machine Learning DataStore.

## Example Usage with Azure Blob

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

resource "azurerm_storage_container" "example" {
  name                  = "example-container"
  storage_account_name  = azurerm_storage_account.example.name
  container_access_type = "private"
}

resource "azurerm_machine_learning_datastore" "example" {
  name                 = "example-datastore"
  resource_group_name  = azurerm_resource_group.example.name
  workspace_name       = azurerm_machine_learning_workspace.example.name
  type                 = "AzureBlob"
  storage_account_name = azurerm_storage_account.example.name
  container_name       = azurerm_storage_container.example.name

  credentials {
    account_key = azurerm_storage_account.example.primary_access_key
  }
}
```

## Example Usage with Azure File

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

resource "azurerm_storage_share" "example" {
  name                 = "example-share"
  storage_account_name = azurerm_storage_account.example.name
  quota                = 1
}

resource "azurerm_machine_learning_datastore" "example" {
  name                 = "example-datastore"
  resource_group_name  = azurerm_resource_group.example.name
  workspace_name       = azurerm_machine_learning_workspace.example.name
  type                 = "AzureFile"
  storage_account_name = azurerm_storage_account.example.name
  file_share_name      = azurerm_storage_share.example.name

  credentials {
    account_key = azurerm_storage_account.example.primary_access_key
  }
}
```

## Example Usage with Data Lake Gen2

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

resource "azurerm_storage_container" "example" {
  name                  = "example-container"
  storage_account_name  = azurerm_storage_account.example.name
  container_access_type = "private"
}

resource "azurerm_machine_learning_datastore" "example" {
  name                 = "example-datastore"
  resource_group_name  = azurerm_resource_group.example.name
  workspace_name       = azurerm_machine_learning_workspace.example.name
  type                 = "AzureDataLakeGen2"
  storage_account_name = azurerm_storage_account.example.name
  container_name       = azurerm_storage_container.example.name
}
```

## Example Usage with Data Lake Gen2 authorize with SPN

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

resource "azurerm_storage_container" "example" {
  name                  = "example-container"
  storage_account_name  = azurerm_storage_account.example.name
  container_access_type = "private"
}

resource "azurerm_machine_learning_datastore" "example" {
  name                 = "example-datastore"
  resource_group_name  = azurerm_resource_group.example.name
  workspace_name       = azurerm_machine_learning_workspace.example.name
  type                 = "AzureDataLakeGen2"
  storage_account_name = azurerm_storage_account.example.name
  container_name       = azurerm_storage_container.example.name

  credentials {
    tenant_id     = "00000000-0000-0000-0000-000000000000"
    client_id     = "00000000-0000-0000-0000-000000000000"
    client_secret = "00000000-0000-0000-0000-000000000000"
  }
}
```
## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Machine Learning DataStore. Changing this forces a new Machine Learning DataStore to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Machine Learning DataStore should exist. Changing this forces a new Machine Learning DataStore to be created.

* `type` - (Required) The type of the data store. Changing this forces a new Machine Learning DataStore to be created. Possible values are `AzureBlob`, `AzureFile`, `AzureDataLakeGen1` and `AzureDataLakeGen2`.

* `workspace_name` - (Required) The name of the machine learning workspace. Changing this forces a new Machine Learning DataStore to be created.

---
* `storage_account_name` - (Optional) The name of the storage account. Changing this forces a new Machine Learning DataStore to be created.

* `container_name` - (Optional) The name of the storage account container. Changing this forces a new Machine Learning DataStore to be created.

* `credentials` - (Optional) A `credentials` block as defined below.

* `file_share_name` - (Optional) The name of the storage account file share. Changing this forces a new Machine Learning DataStore to be created.

* `description` - (Optional) Text used to describe the asset. Changing this forces a new Machine Learning DataStore to be created.

* `is_default` - (Optional) A bool indicate if datastore is the workspace default datastore. Defaults to `false`.

~> **Note:** `is_default` can only be set to `true` on update. 

* `service_data_auth_identity` - (Optional) Indicates which identity to use to authenticate service data access to customer's storage. Defaults to `None`. Possible values are `None`, `WorkspaceSystemAssignedIdentity` and `WorkspaceUserAssignedIdentity`.

* `store_name` - (Optional) The name of Data Lake Storage Gen1. Changing this forces a new Machine Learning DataStore to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Machine Learning DataStore. Changing this forces a new Machine Learning DataStore to be created.

---

A `credentials` block supports the following:

* `account_key` - (Optional) The access key of the storage account. Supported by data store type `AzureBlob` and `AzureFile`. Conflicts with `shared_access_signature`,`tenant_id`,`client_id`,`client_secret`

* `shared_access_signature` - (Optional) The shared access signature of the storage account. Supported by data store type `AzureBlob` and `AzureFile`. Conflicts with `account_key`,`tenant_id`,`client_id`,`client_secret`

* `tenant_id` - (Optional) The ID of the tenant. Supported by data store type `AzureDataLakeGen1` and `AzureDataLakeGen2`. Conflicts with `account_key`,`shared_access_signature`

* `client_id` - (Optional) The ID of the service principal. Supported by data store type `AzureDataLakeGen1` and `AzureDataLakeGen2`. Conflicts with `account_key`,`shared_access_signature`

* `client_secret` - (Optional) The secret of the service principal. Supported by data store type `AzureDataLakeGen1` and `AzureDataLakeGen2`. Conflicts with `account_key`,`shared_access_signature`

* `authority_url` - (Optional) An URL used for authentication. Supported by data store type `AzureDataLakeGen1` and `AzureDataLakeGen2`. Conflicts with `account_key`,`shared_access_signature`.

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
terraform import azurerm_machine_learning_datastore.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.MachineLearningServices/workspaces/mlw1/datastores/datastore1
```
