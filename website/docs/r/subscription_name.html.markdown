---
subcategory: "Base"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_subscription_name"
description: |-
  Manages a Subscription's name and tags.
---

# azurerm_subscription_name

Manages a Subscription's name and tags, by adding an Alias to an existing Subscription.

When this resource is destroyed, it will delete the Alias.
It does not cancel the subscription, revert the name or remove tags from the subscription.

If you want to create a subscription, use the `azurerm_subscription` resource.

~> **NOTE:** This resource **cannot** be used when a subscription is managed using `azurerm_subscription`, as both will attempt to manage the subscription.

~> **NOTE:** Azure supports Multiple Aliases per Subscription, however, to reliably manage this resource in Terraform only a single Alias is supported.

## Example Usage - renaming an existing Subscription

```hcl
resource "azurerm_subscription_name" "example" {
  subscription_id = "12345678-12234-5678-9012-123456789012"
  name            = "My Example Subscription"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The Name of the Subscription. This is the Display Name in the portal.

* `subscription_id` - (Required) The ID of the Subscription. Changing this forces a new Subscription Alias to be created.

---

* `alias` - (Optional) The Alias name for the subscription. Terraform will generate a new GUID if this is not supplied. Changing this forces a new Subscription Alias to be created.

* `tags` - (Optional) A mapping of tags to assign to the Subscription.

* `workload` - (Optional) The workload type of the Subscription.  Possible values are `Production` (default) and `DevTest`. Changing this forces a new Subscription Alias to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The Resource ID of the Alias.

* `tenant_id` - The ID of the Tenant to which the subscription belongs.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Subscription.
* `read` - (Defaults to 5 minutes) Used when retrieving the Subscription.
* `update` - (Defaults to 30 minutes) Used when updating the Subscription.
* `delete` - (Defaults to 30 minutes) Used when deleting the Subscription.

## Import

Subscriptions can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_subscription_name.example "/providers/Microsoft.Subscription/aliases/subscription1"
```

!> **NOTE:** When importing a Subscription that was not created programmatically, it will have no Alias ID to import via `terraform import`.
In this scenario, you can define a resource and Terraform will assume control of the existing subscription by creating an Alias.
