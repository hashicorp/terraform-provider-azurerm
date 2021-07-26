---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_local_network_gateway"
description: |-
  Gets information about an existing Local Network Gateway.
---

# Data Source: azurerm_local_network_gateway

Use this data source to access information about an existing Local Network Gateway.

## Example Usage

```hcl
data "azurerm_local_network_gateway" "example" {
  name                = "existing-local-network-gateway"
  resource_group_name = "existing-resources"
}

output "id" {
  value = data.azurerm_local_network_gateway.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Local Network Gateway.

* `resource_group_name` - (Required) The name of the Resource Group where the Local Network Gateway exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Local Network Gateway.

* `location` - The Azure Region where the Local Network Gateway exists.

* `address_space` - The list of string CIDRs representing the address spaces the gateway exposes.

* `bgp_settings` - A `bgp_settings` block as defined below containing the Local Network Gateway's BGP speaker settings.

* `gateway_address` - The gateway IP address the Local Network Gateway uses.

* `gateway_fqdn` - The gateway FQDN the Local Network Gateway uses.

* `tags` - A mapping of tags assigned to the Local Network Gateway.

---

`bgp_settings` exports the following:

* `asn` - The BGP speaker's ASN.

* `bgp_peering_address` - The BGP peering address and BGP identifier of this BGP speaker.

* `peer_weight` - The weight added to routes learned from this BGP speaker.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Local Network Gateway.