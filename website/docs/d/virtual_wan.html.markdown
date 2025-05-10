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
  value = data.azurerm_virtual_wan.example.disable_vpn_encryption
}

output "location" {
  value = data.azurerm_virtual_wan.example.location
}

output "office365_local_breakout_category" {
  value = data.azurerm_virtual_wan.example.office365_local_breakout_category
}

output "sku" {
  value = data.azurerm_virtual_wan.example.sku
}

output "tags" {
  value = data.azurerm_virtual_wan.example.tags
}

output "virtual_hubs" {
  value = data.azurerm_virtual_wan.example.virtual_hubs
}

output "vpn_sites" {
  value = data.azurerm_virtual_wan.example.vpn_sites
}

```

## Arguments Reference

The following arguments are supported:

- `name` - (Required) The name of this Virtual Wan.

- `resource_group_name` - (Required) The name of the Resource Group where the Virtual Wan exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

- `id` - The ID of the Virtual Wan.

- `allow_branch_to_branch_traffic` - Is branch to branch traffic is allowed?

- `disable_vpn_encryption` - Is VPN Encryption disabled?

- `location` - The Azure Region where the Virtual Wan exists.

- `office365_local_breakout_category` - The Office365 Local Breakout Category.

- `sku` - Type of Virtual Wan (Basic or Standard).

- `tags` - A mapping of tags assigned to the Virtual Wan.

- `virtual_hub_ids` - A list of Virtual Hubs IDs attached to this Virtual WAN.

- `vpn_site_ids` - A list of VPN Site IDs attached to this Virtual WAN.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Virtual Wan.
