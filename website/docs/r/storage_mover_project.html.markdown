---
subcategory: "Storage Mover"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_mover_project"
description: |-
  Manages a Storage Mover Project.
---

# azurerm_storage_mover_project

Manages a Storage Mover Project.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_mover" "example" {
  name                = "example-ssm"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_storage_mover_project" "example" {
  name             = "example-sp"
  storage_mover_id = azurerm_storage_mover.example.id
  description      = "Example Project Description"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Storage Mover Project. Changing this forces a new resource to be created.

* `storage_mover_id` - (Required) Specifies the ID of the storage mover for this Storage Mover Project. Changing this forces a new resource to be created.

* `description` - (Optional) Specifies a description for this Storage Mover Project.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Storage Mover Project.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Storage Mover Project.
* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Mover Project.
* `update` - (Defaults to 30 minutes) Used when updating the Storage Mover Project.
* `delete` - (Defaults to 30 minutes) Used when deleting the Storage Mover Project.

## Import

Storage Mover Project can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_storage_mover_project.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.StorageMover/storageMovers/storageMover1/projects/project1
```
