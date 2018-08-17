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
resource "azurerm_policy_definition" "policy" {
    name = "TestManagementGroup"
    subscription_ids = [
        "000000-1111-2222-3333-444444444444"
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

Management groups can be imported using the `management group name`, e.g.

```shell
terraform import azurerm_management_group.testManagementGroup  /providers/Microsoft.Management/ManagementGroups/<MANAGEMENT_GROUP_NAME>
```
