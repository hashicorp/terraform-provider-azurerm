---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_front_door"
sidebar_current: "docs-azurerm-resource-front-door"
description: |-
  Managed a Front Door on Azure.
---

# azurerm_front_door

Managed a Front Door on Azure.


## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "West US"
}

resource "azurerm_frontdoor" "example" {
  name                = "example-frontdoor"
  resource_group_name = "${azurerm_resource_group.example.name}"
  location            = "${azurerm_resource_group.example.location}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the Front Door which is globally unique. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Batch Account. Changing this forces a new resource to be created.

* `location` - (Required) Resource location. Changing this forces a new resource to be created.

* `enabled_state` - (Optional) Operational status of the Front Door load balancer. Defaults to `Enabled`.

* `enforce_certificate_name_check` - (Optional) Whether to enforce certificate name check on HTTPS requests to all backend pools. No effect on non-HTTPS requests.

* `friendly_name` - (Optional) A friendly name for the frontDoor.

* `tags` - (Optional) Resource tags.

## Attributes Reference

The following attributes are exported:

* `id` - Resource ID.


## Import

Front Door can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_front_door.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.Network/frontDoors/example-frontdoor
```
