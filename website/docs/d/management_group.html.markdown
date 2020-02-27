---
subcategory: "Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_management_group"
description: |-
  Gets information about an existing Management Group.
---

# Data Source: azurerm_management_group

Use this data source to access information about an existing Management Group.

## Example Usage

```hcl
data "azurerm_management_group" "example" {
  group_id = "00000000-0000-0000-0000-000000000000"
}

output "display_name" {
  value = data.azurerm_management_group.example.display_name
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - Specifies the UUID of this Management Group.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Management Group.

* `display_name` - A friendly name for the Management Group.

* `parent_management_group_id` - The ID of any Parent Management Group.

* `subscription_ids` - A list of Subscription ID's which are assigned to the Management Group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Management Group.
