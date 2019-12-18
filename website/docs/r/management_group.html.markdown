---
subcategory: "Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_management_group"
sidebar_current: "docs-azurerm-management-group"
description: |-
  Manages a Management Group.
---

# azurerm_management_group

Manages a Management Group.

## Example Usage

```hcl
data "azurerm_subscription" "current" {}

resource "azurerm_management_group" "example_parent" {
  display_name = "ParentGroup"

  subscription_ids = [
    "${data.azurerm_subscription.current.subscription_id}",
  ]
}

resource "azurerm_management_group" "example_child" {
  display_name               = "ChildGroup"
  parent_management_group_id = "${azurerm_management_group.example_parent.id}"

  subscription_ids = [
    "${data.azurerm_subscription.current.subscription_id}",
  ]

  # other subscription IDs can go here
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Optional) The UUID for this Management Group, which needs to be unique across your tenant - which will be generated if not provided. Changing this forces a new resource to be created.

* `display_name` - (Optional) A friendly name for this Management Group. If not specified, this'll be the same as the `group_id`.

* `parent_management_group_id` - (Optional) The ID of the Parent Management Group. Changing this forces a new resource to be created.

* `subscription_ids` - (Optional) A list of Subscription GUIDs which should be assigned to the Management Group.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Management Group.

## Import

Management Groups can be imported using the `management group resource id`, e.g.

```shell
terraform import azurerm_management_group.example /providers/Microsoft.Management/ManagementGroups/group1
```
