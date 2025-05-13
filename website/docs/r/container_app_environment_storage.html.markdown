---
subcategory: "Container Apps"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_container_app_environment_storage"
description: |-
  Manages a Container App Environment Storage.
---

# azurerm_container_app_environment_storage

Manages a Container App Environment Storage.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_log_analytics_workspace" "example" {
  name                = "acctest-01"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_container_app_environment" "example" {
  name                       = "myEnvironment"
  location                   = azurerm_resource_group.example.location
  resource_group_name        = azurerm_resource_group.example.name
  log_analytics_workspace_id = azurerm_log_analytics_workspace.example.id
}

resource "azurerm_storage_account" "example" {
  name                     = "azureteststorage"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_share" "example" {
  name                 = "sharename"
  storage_account_name = azurerm_storage_account.example.name
  quota                = 5
}

resource "azurerm_container_app_environment_storage" "example" {
  name                         = "mycontainerappstorage"
  container_app_environment_id = azurerm_container_app_environment.example.id
  account_name                 = azurerm_storage_account.example.name
  share_name                   = azurerm_storage_share.example.name
  access_key                   = azurerm_storage_account.example.primary_access_key
  access_mode                  = "ReadOnly"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name for this Container App Environment Storage. Changing this forces a new resource to be created.

* `container_app_environment_id` - (Required) The ID of the Container App Environment to which this storage belongs. Changing this forces a new resource to be created.

* `account_name` - (Optional) The Azure Storage Account in which the Share to be used is located. Changing this forces a new resource to be created.

* `access_key` - (Optional) The Storage Account Access Key.

* `share_name` - (Required) The name of the Azure Storage Share to use. Changing this forces a new resource to be created.

* `access_mode` - (Required) The access mode to connect this storage to the Container App. Possible values include `ReadOnly` and `ReadWrite`. Changing this forces a new resource to be created.

* `nfs_server_url` - (Optional) The NFS server to use for the Azure File Share, the format will be `yourstorageaccountname.file.core.windows.net`. Changing this forces a new resource to be created.
* 
## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Container App Environment Storage


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Container App Environment Storage.
* `read` - (Defaults to 5 minutes) Used when retrieving the Container App Environment Storage.
* `update` - (Defaults to 30 minutes) Used when updating the Container App Environment Storage.
* `delete` - (Defaults to 30 minutes) Used when deleting the Container App Environment Storage.

## Import

A Container App Environment Storage can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_container_app_environment_storage.example "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.App/managedEnvironments/myEnvironment/storages/mystorage"
```
