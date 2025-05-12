---
subcategory: "Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_management_group"
description: |-
  Manages a Management Group.
---

# azurerm_management_group

Manages a Management Group.

!> **Note:** Configuring `subscription_ids` is not supported when using the `azurerm_management_group_subscription_association` resource, results will be unpredictable.

## Example Usage

```hcl
data "azurerm_subscription" "current" {
}

resource "azurerm_management_group" "example_parent" {
  display_name = "ParentGroup"

  subscription_ids = [
    data.azurerm_subscription.current.subscription_id,
  ]
}

resource "azurerm_management_group" "example_child" {
  display_name               = "ChildGroup"
  parent_management_group_id = azurerm_management_group.example_parent.id

  subscription_ids = [
    data.azurerm_subscription.current.subscription_id,
  ]
  # other subscription IDs can go here
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name or UUID for this Management Group, which needs to be unique across your tenant. A new UUID will be generated if not provided. Changing this forces a new resource to be created.

* `display_name` - (Optional) A friendly name for this Management Group. If not specified, this will be the same as the `name`.

* `parent_management_group_id` - (Optional) The ID of the Parent Management Group.

* `subscription_ids` - (Optional) A list of Subscription GUIDs which should be assigned to the Management Group.

~> **Note:** To clear all Subscriptions from the Management Group set `subscription_ids` to an empty list

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Management Group.

* `tenant_scoped_id` - The Management Group ID with the Tenant ID prefix.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Management Group.
* `read` - (Defaults to 5 minutes) Used when retrieving the Management Group.
* `update` - (Defaults to 30 minutes) Used when updating the Management Group.
* `delete` - (Defaults to 30 minutes) Used when deleting the Management Group.

## Import

Management Groups can be imported using the `management group resource id`, e.g.

```shell
terraform import azurerm_management_group.example /providers/Microsoft.Management/managementGroups/group1
```
