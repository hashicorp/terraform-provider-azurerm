---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_hub"
description: |-
  Manages a Virtual Hub within a Virtual WAN.
---

# azurerm_virtual_hub

Manages a Virtual Hub within a Virtual WAN.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_wan" "example" {
  name                = "example-virtualwan"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_virtual_hub" "example" {
  name                = "example-virtualhub"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  virtual_wan_id      = azurerm_virtual_wan.example.id
  address_prefix      = "10.0.0.0/23"
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) The name of the Virtual Hub. Changing this forces a new resource to be created.

- `resource_group_name` - (Required) Specifies the name of the Resource Group where the Virtual Hub should exist. Changing this forces a new resource to be created.

- `location` - (Required) Specifies the supported Azure location where the Virtual Hub should exist. Changing this forces a new resource to be created.

---

- `address_prefix` - (Optional) The Address Prefix which should be used for this Virtual Hub. Changing this forces a new resource to be created. [The address prefix subnet cannot be smaller than a `/24`. Azure recommends using a `/23`](https://docs.microsoft.com/en-us/azure/virtual-wan/virtual-wan-faq#what-is-the-recommended-hub-address-space-during-hub-creation).

- `default_route_table` - (Optional) A`default_route_table` block as defined below.

- `route` - (Optional) One or more `route` blocks as defined below.

- `sku` - (Optional) The sku of the Virtual Hub. Possible values are `Basic` and `Standard`. Changing this forces a new resource to be created.

- `virtual_wan_id` - (Optional) The ID of a Virtual WAN within which the Virtual Hub should be created. Changing this forces a new resource to be created.

- `tags` - (Optional) A mapping of tags to assign to the Virtual Hub.

The `default_route_table` block supports the following:

- `labels` - (Optional) List of labels associated with the default route table.

- `route` - (Optional) One or more `route` blocks as defined below.

---

A `route` block within the `default_route_table` block supports the following:

- `name` - (Required) The name which should be used for this route.

- `destinations` - (Required) A list of destination addresses for this route.

- `destinations_type` - (Required) The type of destinations. Possible values are `CIDR`, `ResourceId` and `Service`.

- `next_hop` - (Required) The next hop's resource ID.

- `next_hop_type` - (Optional) The type of next hop. Currently the only possible value is `ResourceId`. Defaults to `ResourceId`.

---

The `route` block supports the following:

- `address_prefixes` - (Required) A list of Address Prefixes.

- `next_hop_ip_address` - (Required) The IP Address that Packets should be forwarded to as the Next Hop.

---

## Attributes Reference

The following attributes are exported:

- `id` - The ID of the Virtual Hub.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

- `create` - (Defaults to 60 minutes) Used when creating the Virtual Hub.
- `update` - (Defaults to 60 minutes) Used when updating the Virtual Hub.
- `read` - (Defaults to 5 minutes) Used when retrieving the Virtual Hub.
- `delete` - (Defaults to 60 minutes) Used when deleting the Virtual Hub.

## Import

Virtual Hub's can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_virtual_hub.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/virtualHubs/hub1
```
