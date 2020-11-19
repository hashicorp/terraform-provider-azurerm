---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_virtual_wan"
description: |-
  Gets information about an existing Virtual Wan.
---

# Data Source: azurerm_virtual_wan

Use this data source to access information about an existing Virtual Wan.

## Example Usage

```hcl
data "azurerm_virtual_wan" "example" {
  name                = "existing"
  resource_group_name = "existing"
}

output "id" {
  value = data.azurerm_virtual_wan.example.id
}

output "allow_branch_to_branch_traffic" {
  value = data.azurerm_virtual_wan.example.allow_branch_to_branch_traffic
}

output "disable_vpn_encryption" {
  value = data.azurerm_virtual_wan.exemple.disable_vpn_encryption
}

output "location" {
  value = data.azurerm_virtual_wan.exemple.location
}

output "office365_local_breakout_category" {
  value = data.azurerm_virtual_wan.exemple.office365_local_breakout_category
}

output "sku" {
  value = data.azurerm_virtual_wan.exemple.sku
}

output "tags" {
  value = data.azurerm_virtual_wan.exemple.tags
}

output "virtual_hubs" {
  value = data.azurerm_virtual_wan.exemple.virtual_hubs
}

output "vpn_sites" {
  value = data.azurerm_virtual_wan.exemple.vpn_sites
}

```

## Arguments Reference

The following arguments are supported:

- `name` - (Required) The name of this Virtual Wan.

- `resource_group_name` - (Required) The name of the Resource Group where the Virtual Wan exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

- `id` - The ID of the Virtual Wan.

- `allow_branch_to_branch_traffic` - (Optional) True, if allow branch to branch traffic is allowed.

- `disable_vpn_encryption` - (Optional) True, if vpn encryption is disabled.

- `location` - The Azure Region where the Virtual Wan exists.

- `office365_local_breakout_category` - (Optional) Possible value are :
    - `All`
    - `None`
    - `Optimize`
    - `OptimizeAndAllow`.

- `sku` - Type of Virtual Wan (Basic or Standard).

- `tags` - A mapping of tags assigned to the Virtual Wan.

- `virtual_hub_ids` - A list of Virtual Hubs ID's attached to this Virtual WAN.

- `vpn_site_ids` - A list of VPN Site ID's attached to this Virtual WAN.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

- `read` - (Defaults to 5 minutes) Used when retrieving the Virtual Wan.
