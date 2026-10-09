---
subcategory: "Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_management_lock"
description: |-
  Lists Management Lock resources.
---

# List resource: azurerm_management_lock

Lists Management Lock resources.

## Example Usage

### List all Management Locks applied at the Subscription level

-> **Note:** Listing at the Subscription level queries locks applied directly to the subscription itself, rather than searching child resources or resource groups across the subscription.

```hcl
list "azurerm_management_lock" "example" {
  provider = azurerm
  config {
    subscription_id = "00000000-0000-0000-0000-000000000000"
  }
}
```

### List all Management Locks applied at a Resource Group level

```hcl
list "azurerm_management_lock" "example" {
  provider = azurerm
  config {
    resource_group_name = "group1"
  }
}
```

### List all Management Locks applied to a Resource

```hcl
list "azurerm_management_lock" "example" {
  provider = azurerm
  config {
    resource_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/virtualMachines/vm1"
  }
}
```

### List all Management Locks on a Scope

```hcl
list "azurerm_management_lock" "example" {
  provider = azurerm
  config {
    scope = "/providers/Microsoft.Management/managementGroups/group1"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `subscription_id` - (Optional) The ID of the Subscription to query locks for (supports either a raw UUID or `/subscriptions/{subscriptionId}`). Queries locks applied directly at the subscription level. Conflicts with `resource_group_name`, `resource_id`, and `scope`. If all arguments are omitted, defaults to querying the current subscription at the subscription level.

* `resource_group_name` - (Optional) The name of the Resource Group to query locks for. Queries locks applied directly at the resource group level. Conflicts with `subscription_id`, `resource_id`, and `scope`.

* `resource_id` - (Optional) The ID of the Resource to query locks for. Queries locks applied directly at the resource level. Conflicts with `subscription_id`, `resource_group_name`, and `scope`.

* `scope` - (Optional) The Scope to query locks for (e.g. a Management Group or arbitrary resource scope). Conflicts with `subscription_id`, `resource_group_name`, and `resource_id`.
