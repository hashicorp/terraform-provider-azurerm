---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_container"
description: |-
  Manages a Container within an Azure Storage Account.
---

# azurerm_storage_container

Manages a Container within an Azure Storage Account.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "examplestoraccount"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_container" "example" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.example.name
  container_access_type = "private"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Container which should be created within the Storage Account.

* `storage_account_name` - (Required) The name of the Storage Account where the Container should be created.

* `container_access_type` - (Optional) The Access Level configured for this Container. Possible values are `blob`, `container` or `private`. Defaults to `private`.

* `metadata` - (Optional) A mapping of MetaData for this Container.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the Storage Container.

* `has_immutability_policy` - Is there an Immutability Policy configured on this Storage Container?

* `has_legal_hold` - Is there a Legal Hold configured on this Storage Container?

* `resource_manager_id` - The Resource Manager ID of this Storage Container.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Storage Container.
* `update` - (Defaults to 30 minutes) Used when updating the Storage Container.
* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Container.
* `delete` - (Defaults to 30 minutes) Used when deleting the Storage Container.

## Import

Storage Containers can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_storage_container.container1 https://example.blob.core.windows.net/container
```
