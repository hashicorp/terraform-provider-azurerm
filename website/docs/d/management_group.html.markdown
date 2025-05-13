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
  name = "00000000-0000-0000-0000-000000000000"
}

output "display_name" {
  value = data.azurerm_management_group.example.display_name
}
```

## Argument Reference

The following arguments are supported:

* `name` - Specifies the name or UUID of this Management Group.

* `display_name` - Specifies the display name of this Management Group.

~> **Note:** Whilst multiple management groups may share the same display name, when filtering Terraform expects a single management group to be found with this name.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Management Group.

* `tenant_scoped_id` - The Management Group ID with the Tenant ID prefix.

* `parent_management_group_id` - The ID of any Parent Management Group.

* `management_group_ids` - A list of Management Group IDs which directly belong to this Management Group.

* `subscription_ids` - A list of Subscription IDs which are directly assigned to this Management Group.

* `all_management_group_ids` - A list of Management Group IDs which directly or indirectly belong to this Management Group.

* `all_subscription_ids` - A list of Subscription IDs which are assigned to this Management Group or its children Management Groups.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Management Group.
