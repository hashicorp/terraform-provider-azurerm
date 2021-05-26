---
subcategory: "Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_management_group_subscription_association"
description: |-
  Manages a Management Group Subscription Association.
---

# azurerm_management_group_subscription_association

Manages a Management Group Subscription Association.

!> **Note:** When using this resource, configuring `subscription_ids` on the `azurerm_management_group` resource is not supported.

## Example Usage

```hcl
data "azurerm_management_group" "example" {
  name = "exampleManagementGroup"
}

data "azurerm_subscription" "example" {
  subscription_id = "12345678-1234-1234-1234-123456789012"
}

resource "azurerm_management_group_subscription_association" "example" {
  management_group_id = data.azurerm_management_group.example.id
  subscription_id     = data.azurerm_subscription.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `management_group_id` - (Required) The ID of the Management Group to associate the Subscription with. Changing this forces a new Management to be created.

* `subscription_id` - (Required) The ID of the Subscription to be associated with the Management Group. Changing this forces a new Management to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Management Group Subscription Association.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the Management.
* `read` - (Defaults to 5 minutes) Used when retrieving the Management.
* `delete` - (Defaults to 5 minutes) Used when deleting the Management.

## Import

Managements can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_management_group_subscription_association.example /managementGroup/MyManagementGroup/subscription/12345678-1234-1234-1234-123456789012
```
