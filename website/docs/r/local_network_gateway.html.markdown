---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_local_network_gateway"
description: |-
  Manages a local network gateway connection over which specific connections can be configured.
---

# azurerm_local_network_gateway

Manages a local network gateway connection over which specific connections can be configured.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "localNetworkGWTest"
  location = "West US"
}

resource "azurerm_local_network_gateway" "home" {
  name                = "backHome"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  gateway_address     = "12.13.14.15"
  address_space       = ["10.0.0.0/16"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the local network gateway. Changing this
    forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to
    create the local network gateway.

* `location` - (Required) The location/region where the local network gateway is
    created. Changing this forces a new resource to be created.

* `address_space` - (Required) The list of string CIDRs representing the
    address spaces the gateway exposes.

* `bgp_settings` - (Optional) A `bgp_settings` block as defined below containing the
    Local Network Gateway's BGP speaker settings.
    
* `gateway_address` - (Optional) The gateway IP address to connect with.
    
* `gateway_fqdn` - (Optional) The gateway FQDN to connect with.

-> **NOTE**: Either `gateway_address` or `gateway_fqdn` should be specified.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

`bgp_settings` supports the following:

* `asn` - (Required) The BGP speaker's ASN.

* `bgp_peering_address` - (Required) The BGP peering address and BGP identifier
    of this BGP speaker.

* `peer_weight` - (Optional) The weight added to routes learned from this
    BGP speaker.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Local Network Gateway.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Local Network Gateway.
* `update` - (Defaults to 30 minutes) Used when updating the Local Network Gateway.
* `read` - (Defaults to 5 minutes) Used when retrieving the Local Network Gateway.
* `delete` - (Defaults to 30 minutes) Used when deleting the Local Network Gateway.

## Import

Local Network Gateways can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_local_network_gateway.lng1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/localNetworkGateways/lng1
```
