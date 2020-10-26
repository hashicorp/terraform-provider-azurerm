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

  vpn_site_link {
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

* `vpn_site_link` - (Required) One or more `vpn_site_link` blocks as defined below.

---

* `address_spaces` - (Optional) Specifies a list of IP address spaces that is located on your on-premises site. Traffic destined for the address spaces is routed to your local site.

* `device_property` - (Optional) A `device_property` block as defined below.

* `tags` - (Optional) A mapping of tags which should be assigned to the VPN Site.

---

A `bgp_property` block supports the following:

* `asn` - (Required) The BGP speaker's ASN.

* `peering_address` - (Required) The BGP peering ip address. 

---

A `device_property` block supports the following:

* `model` - (Optional) The model of the VPN device.

* `vendor` - (Optional) The name of the VPN device vendor.
                        
---

A `vpn_site_link` block supports the following:

* `name` - (Required) The name which should be used for this VPN Site Link.

* `bgp_property` - (Optional) A `bgp_property` block as defined above.

* `fqdn` - (Optional) The FQDN of this VPN Site Link.

* `ip_address` - (Optional) The IP address of this VPN Site Link.

-> **NOTE**: Either `fqdn` or `ip_address` should be specified.

* `link_provider_name` - (Optional) The name of the physical link at the VPN Site. Example: ATT, Verizon.

* `link_speed_mbps` - (Optional) The speed of the VPN device at the branch location in unit of mbps.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the VPN Site.

* `vpn_site_link` - One or more `vpn_site_link` blocks as defined below.

---

A `vpn_site_link` block supports the following:

* `id` - The ID of the VPN Site Link.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the VPN Site.
* `read` - (Defaults to 5 minutes) Used when retrieving the VPN Site.
* `update` - (Defaults to 30 minutes) Used when updating the VPN Site.
* `delete` - (Defaults to 30 minutes) Used when deleting the VPN Site.

## Import

VPN Sites can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_vpn_site.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/vpnSites/site1
```
