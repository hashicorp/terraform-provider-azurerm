---
subcategory: "PowerBIDedicated"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_powerbidedicated_capacity"
sidebar_current: "docs-azurerm-resource-powerbidedicated-capacity"
description: |-
  Manages a PowerBIDedicated Capacity.
---

# azurerm_powerbidedicated_capacity

Manages a PowerBIDedicated Capacity.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_powerbidedicated_capacity" "example" {
  name                = "example-powerbidedicatedcapacity"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  sku                 = "A1"
  administrators      = ["azsdktest@microsoft.com"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the PowerBIDedicated Capacity. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group where the PowerBIDedicated Capacity should be created. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `sku` - (Required) The SKU of the PowerBIDedicated Capacity. Valid values include `A1`, `A2`, `A3`, `A4`, `A5` or `A6`.

* `administrators` - (Required) A set of administrator user identities.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the PowerBIDedicated Capacity.

## Import

PowerBIDedicated Capacities can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_powerbidedicated_capacity.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.PowerBIDedicated/capacities/capacity1
```
