---
subcategory: "Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_management_group_subscription_association"
description: |-
  Associates a Subscription ID to a Management Group
---

# azurerm_management_group_subscription_association

Associates a Subscription ID to a Management Group.

## Example Usage

```hcl
data "azurerm_subscription" "current" {}

resource "azurerm_management_group" "example" {
  name = "example"

  lifecycle {
    ignore_changes = [subscription_ids, ]
  }
}

resource "azurerm_management_group_subscription_association" "example" {
  management_group_id = azurerm_management_group.example.id
  subscription_id     = data.azurerm_subscription.current.subscription_id
}
```

~> **NOTE:** When using azurerm_management_group_subscription_association you **must** use `lifecycle.ignore_changes = [subscription_ids,]` on the associated azurerm_management_group to avoid conflicts and unexpected plan outputs.

## Arguments Reference

The following arguments are supported:

* `management_group_id` - (Required) The ID of the Management Group. Changing this forces a new resource to be created.

* `subscription_id` - (Required) The subscription GUID to associate to the Management Group. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The (Terraform specific) ID of the Association between the Management Group and the Subscription.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Management Group Subscription Association.
* `read` - (Defaults to 5 minutes) Used when retrieving the Management Group Subscription Association.
* `delete` - (Defaults to 30 minutes) Used when deleting the Management Group Subscription Association.

## Import

Associations between Management Groups and Subscriptions can be imported using the association `resource id`, e.g.

```shell
terraform import azurerm_management_group_subscription_association.example /providers/Microsoft.Management/managementGroups/example|/subscriptions/00000000-0000-0000-0000-000000000000
```

-> **NOTE:** This ID is specific to Terraform - and is of the format `{managementGroupId}|{subscriptionId}`.
