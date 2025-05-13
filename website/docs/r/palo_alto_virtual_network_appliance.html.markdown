---
subcategory: "Palo Alto"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_palo_alto_virtual_network_appliance"
description: |-
  Manages a Palo Alto Network Virtual Appliance.
---

# azurerm_palo_alto_virtual_network_appliance

Manages a Palo Alto Network Virtual Appliance.

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

  tags = {
    "hubSaaSPreview" = "true"
  }
}

resource "azurerm_palo_alto_virtual_network_appliance" "example" {
  name           = "example-appliance"
  virtual_hub_id = azurerm_virtual_hub.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Palo Alto Local Network Virtual Appliance. Changing this forces a new Palo Alto Local Network Virtual Appliance to be created.

* `virtual_hub_id` - (Required) The ID of the Virtual Hub to deploy this appliance onto. Changing this forces a new Palo Alto Local Network Virtual Appliance to be created.

~> **Note:** THe Virtual Hub must be created with the tag `"hubSaaSPreview" = "true"` to be compatible with this resource.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Palo Alto Local Network Virtual Appliance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Palo Alto Local Network Virtual Appliance.
* `read` - (Defaults to 5 minutes) Used when retrieving the Palo Alto Local Network Virtual Appliance.
* `delete` - (Defaults to 30 minutes) Used when deleting the Palo Alto Local Network Virtual Appliance.

## Import

Palo Alto Local Network Virtual Appliances can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_palo_alto_virtual_network_appliance.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/networkVirtualAppliances/myPANetworkVirtualAppliance
```
