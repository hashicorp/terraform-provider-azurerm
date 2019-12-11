---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_hub"
sidebar_current: "docs-azurerm-resource-virtual-hub"
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
  address_prefix      = "10.0.1.0/24"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Virtual Hub. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the Virtual Hub should exist. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the Virtual Hub should exist. Changing this forces a new resource to be created.

* `address_prefix` - (Required) The Address Prefix which should be used for this Virtual Hub.

* `virtual_wan_id` - (Required) The ID of a Virtual WAN within which the Virtual Hub should be created.

---

* `route` - (Optional) One or more `route` blocks as defined below.

* `tags` - (Optional) A mapping of tags to assign to the Virtual Hub.

---

The `route` block supports the following:

* `address_prefixes` - (Required) A list of Address Prefixes.

* `next_hop_ip_address` - (Required) The IP Address that Packets should be forwarded to as the Next Hop.

---

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Virtual Hub.

## Import

Virtual Hub's can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_virtual_hub.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/virtualHubs/hub1
```
