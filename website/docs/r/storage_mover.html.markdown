---
subcategory: "Storage Mover"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_mover"
description: |-
  Manages a Storage Mover.
---

# azurerm_storage_mover

Manages a Storage Mover.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_mover" "example" {
  name                = "example-ssm"
  resource_group_name = azurerm_resource_group.example.name
  location            = "West Europe"
  description         = "Example Storage Mover Description"
  tags = {
    key = "value"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Storage Mover. Changing this forces a new Storage Mover to be created.

* `location` - (Required) Specifies the Azure Region where the Storage Mover should exist. Changing this forces a new Storage Mover to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the Storage Mover should exist. Changing this forces a new Storage Mover to be created.

* `description` - (Optional) A description for the Storage Mover.

* `tags` - (Optional) A mapping of tags which should be assigned to the Storage Mover.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Storage Mover.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Storage Mover.
* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Mover.
* `update` - (Defaults to 30 minutes) Used when updating the Storage Mover.
* `delete` - (Defaults to 30 minutes) Used when deleting the Storage Mover.

## Import

Storage Mover can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_storage_mover.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.StorageMover/storageMovers/storageMover1
```
