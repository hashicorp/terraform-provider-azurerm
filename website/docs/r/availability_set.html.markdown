---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_availability_set"
description: |-
  Manages an Availability Set for Virtual Machines.

---

# azurerm_availability_set

Manages an Availability Set for Virtual Machines.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_availability_set" "example" {
  name                = "example-aset"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  tags = {
    environment = "Production"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the availability set. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the availability set. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `platform_update_domain_count` - (Optional) Specifies the number of update domains that are used. Defaults to `5`. Changing this forces a new resource to be created.

~> **Note:** The number of Update Domains varies depending on which Azure Region you're using. More information about update and fault domains and how they work can be found [here](https://learn.microsoft.com/en-us/azure/virtual-machines/availability-set-overview).

* `platform_fault_domain_count` - (Optional) Specifies the number of fault domains that are used. Defaults to `3`. Changing this forces a new resource to be created.

~> **Note:** The number of Fault Domains varies depending on which Azure Region you're using. More information about update and fault domains and how they work can be found [here](https://learn.microsoft.com/en-us/azure/virtual-machines/availability-set-overview).

* `proximity_placement_group_id` - (Optional) The ID of the Proximity Placement Group to which this Virtual Machine should be assigned. Changing this forces a new resource to be created.

* `managed` - (Optional) Specifies whether the availability set is managed or not. Possible values are `true` (to specify aligned) or `false` (to specify classic). Default is `true`. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Availability Set.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Availability Set.
* `read` - (Defaults to 5 minutes) Used when retrieving the Availability Set.
* `update` - (Defaults to 30 minutes) Used when updating the Availability Set.
* `delete` - (Defaults to 30 minutes) Used when deleting the Availability Set.

## Import

Availability Sets can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_availability_set.group1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Compute/availabilitySets/webAvailSet
```
