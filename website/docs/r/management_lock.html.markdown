---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_management_lock"
sidebar_current: "docs-azurerm-resource-management-lock"
description: |-
  Manages a Management Lock which is scoped to a Subscription, Resource Group or Resource.

---

# azurerm_management_lock

Manages a Management Lock which is scoped to a Subscription, Resource Group or Resource.

## Example Usage

Complete examples of how to use the `azurerm_management_lock` resource can be found [in the `./examples/management-locks` folder within the Github Repository](https://github.com/terraform-providers/terraform-provider-azurerm/tree/master/examples/management-locks)


```hcl
resource "azurerm_resource_group" "test" {
  # ...
}

resource "azurerm_management_lock" "example" {
  name       = "resource-group-level"
  scope      = "${azurerm_resource_group.example.id}"
  lock_level = "ReadOnly"
  notes      = "This Resource Group is Read-Only"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Management Lock. Changing this forces a new resource to be created.

* `scope` - (Required) Specifies the scope at which the Management Lock should be created. Changing this forces a new resource to be created.

* `lock_level` - (Required) Specifies the Level to be used for this Lock. Possible values are `CanNotDelete` and `ReadOnly`. Changing this forces a new resource to be created.

~> **Note:** `CanNotDelete` means authorized users are able to read and modify the resources, but not delete. `ReadOnly` means authorized users can only read from a resource, but they can't modify or delete it.

* `note` - (Optional) Specifies some notes about the lock. Maximum of 512 characters. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Management Lock

## Import

Management Locks can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_management_lock.lock1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Authorization/locks/lock1
```
