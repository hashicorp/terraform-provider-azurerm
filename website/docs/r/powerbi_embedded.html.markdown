---
subcategory: "PowerBI"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_powerbi_embedded"
description: |-
  Manages a PowerBI Embedded.
---

# azurerm_powerbi_embedded

Manages a PowerBI Embedded.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_powerbi_embedded" "example" {
  name                = "examplepowerbi"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "A1"
  administrators      = ["azsdktest@microsoft.com"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the PowerBI Embedded. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the PowerBI Embedded should be created. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `sku_name` - (Required) Sets the PowerBI Embedded's pricing level's SKU. Possible values include: `A1`, `A2`, `A3`, `A4`, `A5`, `A6`.

* `administrators` - (Required) A set of administrator user identities, which manages the Power BI Embedded and must be a member user or a service principal in your AAD tenant.

* `mode` - (Optional) Sets the PowerBI Embedded's mode. Possible values include: `Gen1`, `Gen2`. Defaults to `Gen1`. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the PowerBI Embedded.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the PowerBI Embedded instance.
* `update` - (Defaults to 30 minutes) Used when updating the PowerBI Embedded instance.
* `read` - (Defaults to 5 minutes) Used when retrieving the PowerBI Embedded instance.
* `delete` - (Defaults to 30 minutes) Used when deleting the PowerBI Embedded instance.


## Import

PowerBI Embedded can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_powerbi_embedded.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.PowerBIDedicated/capacities/capacity1
```
