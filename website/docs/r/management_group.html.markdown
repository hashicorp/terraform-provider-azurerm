---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_management_group"
sidebar_current: "docs-azurerm-management-group"
description: |-
  Manages a Management Group.
---

# azurerm_management_group

Create a management group with subscription assignments.

## Example Usage

```hcl

data "azurerm_subscription" "current" {}

resource "azurerm_management_group" "testmanagementgroup" {
    name = "TestManagementGroup"
    subscription_ids = [
        "${data.azurerm_subscription.current.id}"
    ]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name & id of the management group. This needs to be unique across your AAD tenant.

* `subscription_ids` - (Optional) List of subscription IDs to be assigned to the management group.

## Attributes Reference

The following attributes are exported:

* `id` - The management group id.

## Import

Management groups can be imported using the `management group resource id`, e.g.

```shell
terraform import azurerm_management_group.test  /providers/Microsoft.Management/ManagementGroups/group1
```
