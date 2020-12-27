---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_disk_access"
description: |-
  Manages a Disk Access.
---

# azurerm_disk_access

Manages a Disk Access.

## Example Usage

```hcl
resource "azurerm_disk_access" "example" {
  name                = "example"
  resource_group_name = "example"
  location            = "West Europe"
}
```

## Arguments Reference

The following arguments are supported:

* `location` - (Required) The Azure Region where the Disk Access should exist. Changing this forces a new Disk to be created.

* `name` - (Required) The name which should be used for this Disk Access. Changing this forces a new Disk Access to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Disk Access should exist. Changing this forces a new Disk Access to be created.

---

* `tags` - (Optional) A mapping of tags which should be assigned to the Disk Access.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Disk Access resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Disk.
* `read` - (Defaults to 5 minutes) Used when retrieving the Disk.
* `update` - (Defaults to 30 minutes) Used when updating the Disk.
* `delete` - (Defaults to 30 minutes) Used when deleting the Disk.

## Import

Disk Access resource can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_disk_access.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/diskAccesses/diskAccess1
```
