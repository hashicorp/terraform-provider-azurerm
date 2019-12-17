---
subcategory: "PowerBIDedicated"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_powerbi_dedicated_capacity"
sidebar_current: "docs-azurerm-resource-powerbi-dedicated-capacity"
description: |-
  Manages a PowerBI Dedicated Capacity.
---

# azurerm_powerbi_dedicated_capacity

Manages a PowerBI Dedicated Capacity.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_powerbi_dedicated_capacity" "example" {
  name                = "example-powerbidedicatedcapacity"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  sku_name            = "A1"
  administrators      = ["azsdktest@microsoft.com"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the PowerBI Dedicated Capacity. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the PowerBI Dedicated Capacity should be created. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `sku_name` - (Required) Sets the PowerBI Dedicated Capacity's pricing level's SKU. Possible values include: `A1`, `A2`, `A3`, `A4`, `A5`, `A6`.

* `administrators` - (Required) A set of administrator user identities, which manages the capacity in Power BI and must be a member user or a service principal in your AAD tenant.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the PowerBI Dedicated Capacity.

## Import

PowerBI Dedicated Capacities can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_powerbi_dedicated_capacity.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.PowerBIDedicated/capacities/capacity1
```
