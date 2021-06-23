---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dedicated_host_group"
description: |-
  Manage a Dedicated Host Group.
---

# azurerm_dedicated_host_group

Manage a Dedicated Host Group.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-rg-compute"
  location = "West Europe"
}

resource "azurerm_dedicated_host_group" "example" {
  name                        = "example-dedicated-host-group"
  resource_group_name         = azurerm_resource_group.example.name
  location                    = azurerm_resource_group.example.location
  platform_fault_domain_count = 1
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Dedicated Host Group. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the name of the resource group the Dedicated Host Group is located in. Changing this forces a new resource to be created.

* `location` - (Required) The Azure location where the Dedicated Host Group exists. Changing this forces a new resource to be created.

* `platform_fault_domain_count` - (Required) The number of fault domains that the Dedicated Host Group spans. Changing this forces a new resource to be created.

* `automatic_placement_enabled` - (Optional) Would virtual machines or virtual machine scale sets be placed automatically on this Dedicated Host Group? Defaults to `false`. Changing this forces a new resource to be created.

* `zones` - (Optional) A list of Availability Zones in which the Dedicated Host Group should be located. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Dedicated Host Group.

## Timeouts



The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Dedicated Host Group.
* `update` - (Defaults to 30 minutes) Used when updating the Dedicated Host Group.
* `read` - (Defaults to 5 minutes) Used when retrieving the Dedicated Host Group.
* `delete` - (Defaults to 30 minutes) Used when deleting the Dedicated Host Group.

## Import

Dedicated Host Group can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_dedicated_host_group.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.Compute/hostGroups/group1
```
