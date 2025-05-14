---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_vpn_site"
description: |-
  Manages a VPN Site.
---

# azurerm_vpn_site

Manages a VPN Site.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "West Europe"
}

resource "azurerm_virtual_wan" "example" {
  name                = "example-vwan"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_vpn_site" "example" {
  name                = "site1"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  virtual_wan_id      = azurerm_virtual_wan.example.id
  address_cidrs       = ["10.0.0.0/24"]

  link {
    name       = "link1"
    ip_address = "10.0.0.1"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `location` - (Required) The Azure Region where the VPN Site should exist. Changing this forces a new VPN Site to be created.

* `name` - (Required) The name which should be used for this VPN Site. Changing this forces a new VPN Site to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the VPN Site should exist. Changing this forces a new VPN Site to be created.

* `virtual_wan_id` - (Required) The ID of the Virtual Wan where this VPN site resides in. Changing this forces a new VPN Site to be created.

* `link` - (Optional) One or more `link` blocks as defined below.

---

* `address_cidrs` - (Optional) Specifies a list of IP address CIDRs that are located on your on-premises site. Traffic destined for these address spaces is routed to your local site.

-> **Note:** The `address_cidrs` has to be set when the `link.bgp` isn't specified.

* `device_model` - (Optional) The model of the VPN device.

* `device_vendor` - (Optional) The name of the VPN device vendor.

* `o365_policy` - (Optional) An `o365_policy` block as defined below.

* `tags` - (Optional) A mapping of tags which should be assigned to the VPN Site.

---

A `bgp` block supports the following:

* `asn` - (Required) The BGP speaker's ASN.

* `peering_address` - (Required) The BGP peering IP address.

---

A `link` block supports the following:

* `name` - (Required) The name which should be used for this VPN Site Link.

* `bgp` - (Optional) A `bgp` block as defined above.

-> **Note:** The `link.bgp` has to be set when the `address_cidrs` isn't specified.

* `fqdn` - (Optional) The FQDN of this VPN Site Link.

* `ip_address` - (Optional) The IP address of this VPN Site Link.

-> **Note:** Either `fqdn` or `ip_address` should be specified.

* `provider_name` - (Optional) The name of the physical link at the VPN Site. Example: `ATT`, `Verizon`.

* `speed_in_mbps` - (Optional) The speed of the VPN device at the branch location in unit of mbps. Defaults to `0`.

---

A `o365_policy` block supports the following:

* `traffic_category` - (Optional) A `traffic_category` block as defined above.

---

A `traffic_category` block supports the following:

* `allow_endpoint_enabled` - (Optional) Is allow endpoint enabled? The `Allow` endpoint is required for connectivity to specific O365 services and features, but are not as sensitive to network performance and latency as other endpoint types. Defaults to `false`.

* `default_endpoint_enabled` - (Optional) Is default endpoint enabled? The `Default` endpoint represents O365 services and dependencies that do not require any optimization, and can be treated by customer networks as normal Internet bound traffic. Defaults to `false`.

* `optimize_endpoint_enabled` - (Optional) Is optimize endpoint enabled? The `Optimize` endpoint is required for connectivity to every O365 service and represents the O365 scenario that is the most sensitive to network performance, latency, and availability. Defaults to `false`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the VPN Site.

* `link` - One or more `link` blocks as defined below.

---

A `link` block supports the following:

* `id` - The ID of the VPN Site Link.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the VPN Site.
* `read` - (Defaults to 5 minutes) Used when retrieving the VPN Site.
* `update` - (Defaults to 30 minutes) Used when updating the VPN Site.
* `delete` - (Defaults to 30 minutes) Used when deleting the VPN Site.

## Import

VPN Sites can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_vpn_site.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/vpnSites/site1
```
