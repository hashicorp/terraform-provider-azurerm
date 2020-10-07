---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_wan"
description: |-
  Manages a Virtual WAN.

---

# azurerm_virtual_wan

Manages a Virtual WAN.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_wan" "example" {
  name                = "example-vwan"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Virtual WAN. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Virtual WAN. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `disable_vpn_encryption` - (Optional) Boolean flag to specify whether VPN encryption is disabled. Defaults to `false`.

* `allow_branch_to_branch_traffic` - (Optional) Boolean flag to specify whether branch to branch traffic is allowed. Defaults to `true`.

* `office365_local_breakout_category` - (Optional) Specifies the Office365 local breakout category. Possible values include: `Optimize`, `OptimizeAndAllow`, `All`, `None`. Defaults to `None`.

* `type` - (Optional) Specifies the Virtual WAN type. Possible Values include: `Basic` and `Standard`. Defaults to `Standard`.

* `tags` - (Optional) A mapping of tags to assign to the Virtual WAN.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Virtual WAN.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Virtual WAN.
* `update` - (Defaults to 30 minutes) Used when updating the Virtual WAN.
* `read` - (Defaults to 5 minutes) Used when retrieving the Virtual WAN.
* `delete` - (Defaults to 30 minutes) Used when deleting the Virtual WAN.

## Import

Virtual WAN can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_virtual_wan.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/virtualWans/testvwan
```
