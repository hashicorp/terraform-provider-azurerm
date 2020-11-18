---
subcategory: "TODO - pick from: Load Balancer|Network"
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
  name = "existing"
  resource_group_name = "existing"
}

output "id" {
  value = data.azurerm_virtual_wan.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Virtual Wan.

* `resource_group_name` - (Required) The name of the Resource Group where the Virtual Wan exists.

---

* `allow_branch_to_branch_traffic` - (Optional) TODO.

* `disable_vpn_encryption` - (Optional) TODO.

* `office365_local_breakout_category` - (Optional) TODO.

* `virtual_hubs` - (Optional) TODO.

* `vpn_sites` - (Optional) TODO.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Virtual Wan.

* `location` - The Azure Region where the Virtual Wan exists.

* `sku` - TODO.

* `tags` - A mapping of tags assigned to the Virtual Wan.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Virtual Wan.