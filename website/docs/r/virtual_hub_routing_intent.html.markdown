---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_hub_routing_intent"
description: |-
  Manages a Virtual Hub Routing Intent.
---

# azurerm_virtual_hub_routing_intent

Manages a Virtual Hub Routing Intent.

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

resource "azurerm_virtual_hub" "example" {
  name                = "example-vhub"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  virtual_wan_id      = azurerm_virtual_wan.example.id
  address_prefix      = "10.0.1.0/24"
}

resource "azurerm_firewall" "example" {
  name                = "example-fw"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "AZFW_Hub"
  sku_tier            = "Standard"

  virtual_hub {
    virtual_hub_id  = azurerm_virtual_hub.example.id
    public_ip_count = 1
  }
}

resource "azurerm_virtual_hub_routing_intent" "example" {
  name           = "example-routingintent"
  virtual_hub_id = azurerm_virtual_hub.example.id

  routing_policy {
    name         = "InternetTrafficPolicy"
    destinations = ["Internet"]
    next_hop     = azurerm_firewall.example.id
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Virtual Hub Routing Intent. Changing this forces a new resource to be created.

* `virtual_hub_id` - (Required) The resource ID of the Virtual Hub. Changing this forces a new resource to be created.

* `routing_policy` - (Required) One or more `routing_policy` blocks as defined below.

---

A `routing_policy` block supports the following:

* `name` - (Required) The unique name for the routing policy.

* `destinations` - (Required) A list of destinations which this routing policy is applicable to. Possible values are `Internet` and `PrivateTraffic`.

* `next_hop` - (Required) The resource ID of the next hop on which this routing policy is applicable to.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Virtual Hub Routing Intent.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Virtual Hub Routing Intent.
* `read` - (Defaults to 5 minutes) Used when retrieving the Virtual Hub Routing Intent.
* `update` - (Defaults to 30 minutes) Used when updating the Virtual Hub Routing Intent.
* `delete` - (Defaults to 30 minutes) Used when deleting the Virtual Hub Routing Intent.

## Import

Virtual Hub Routing Intents can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_virtual_hub_routing_intent.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Network/virtualHubs/virtualHub1/routingIntent/routingIntent1
```
