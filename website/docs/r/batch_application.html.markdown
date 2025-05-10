---
subcategory: "Batch"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_batch_application"
description: |-
  Manages Azure Batch Application instance.
---

# azurerm_batch_application

Manages Azure Batch Application instance.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "examplesa"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_batch_account" "example" {
  name                                = "exampleba"
  resource_group_name                 = azurerm_resource_group.example.name
  location                            = azurerm_resource_group.example.location
  pool_allocation_mode                = "BatchService"
  storage_account_id                  = azurerm_storage_account.example.id
  storage_account_authentication_mode = "StorageKeys"
}

resource "azurerm_batch_application" "example" {
  name                = "example-batch-application"
  resource_group_name = azurerm_resource_group.example.name
  account_name        = azurerm_batch_account.example.name
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the application. This must be unique within the account. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group that contains the Batch account. Changing this forces a new resource to be created.

* `account_name` - (Required) The name of the Batch account. Changing this forces a new resource to be created.

* `allow_updates` - (Optional) A value indicating whether packages within the application may be overwritten using the same version string. Defaults to `true`.

* `default_version` - (Optional) The package to use if a client requests the application but does not specify a version. This property can only be set to the name of an existing package.

* `display_name` - (Optional) The display name for the application.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Batch Application.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Batch Application.
* `read` - (Defaults to 5 minutes) Used when retrieving the Batch Application.
* `update` - (Defaults to 30 minutes) Used when updating the Batch Application.
* `delete` - (Defaults to 30 minutes) Used when deleting the Batch Application.

## Import

Batch Applications can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_batch_application.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.Batch/batchAccounts/exampleba/applications/example-batch-application
```
