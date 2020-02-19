---
subcategory: "Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_management_lock"
description: |-
  Manages a Management Lock which is scoped to a Subscription, Resource Group or Resource.

---

# azurerm_management_lock

Manages a Management Lock which is scoped to a Subscription, Resource Group or Resource.

## Example Usage (Subscription Level Lock)

```hcl
data "azurerm_subscription" "current" {
}

resource "azurerm_management_lock" "subscription-level" {
  name       = "subscription-level"
  scope      = data.azurerm_subscription.current.id
  lock_level = "CanNotDelete"
  notes      = "Items can't be deleted in this subscription!"
}
```

## Example Usage (Resource Group Level Lock)

```hcl
resource "azurerm_resource_group" "example" {
  name     = "locked-resource-group"
  location = "West Europe"
}

resource "azurerm_management_lock" "resource-group-level" {
  name       = "resource-group-level"
  scope      = azurerm_resource_group.example.id
  lock_level = "ReadOnly"
  notes      = "This Resource Group is Read-Only"
}
```

## Example Usage (Resource Level Lock)

```hcl
resource "azurerm_resource_group" "example" {
  name     = "locked-resource-group"
  location = "West Europe"
}

resource "azurerm_public_ip" "example" {
  name                    = "locked-publicip"
  location                = azurerm_resource_group.example.location
  resource_group_name     = azurerm_resource_group.example.name
  allocation_method       = "Static"
  idle_timeout_in_minutes = 30
}

resource "azurerm_management_lock" "public-ip" {
  name       = "resource-ip"
  scope      = azurerm_public_ip.example.id
  lock_level = "CanNotDelete"
  notes      = "Locked because it's needed by a third-party"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Management Lock. Changing this forces a new resource to be created.

* `scope` - (Required) Specifies the scope at which the Management Lock should be created. Changing this forces a new resource to be created.

* `lock_level` - (Required) Specifies the Level to be used for this Lock. Possible values are `CanNotDelete` and `ReadOnly`. Changing this forces a new resource to be created.

~> **Note:** `CanNotDelete` means authorized users are able to read and modify the resources, but not delete. `ReadOnly` means authorized users can only read from a resource, but they can't modify or delete it.

* `notes` - (Optional) Specifies some notes about the lock. Maximum of 512 characters. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Management Lock

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Management Lock.
* `update` - (Defaults to 30 minutes) Used when updating the Management Lock.
* `read` - (Defaults to 5 minutes) Used when retrieving the Management Lock.
* `delete` - (Defaults to 30 minutes) Used when deleting the Management Lock.

## Import

Management Locks can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_management_lock.lock1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Authorization/locks/lock1
```
