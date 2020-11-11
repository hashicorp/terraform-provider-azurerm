---
subcategory: "Base"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_subscription_alias"
description: |-
  Manages a Subscription Alias.
---

# azurerm_subscription_alias

Manages a Subscription Alias.

## Example Usage

```hcl
data "azurerm_subscription" "current" {}

resource "azurerm_subscription_alias" "example" {
  name            = "example-alias"
  subscription_id = data.azurerm_subscription.current.subscription_id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Subscription Alias. Changing this forces a new resource to be created.

* `subscription_id` - (Optional) The Subscription ID which this Subscription Alias should associate to.

## Attributes Reference

In addition to the arguments listed above - the following attributes are exported: 

* `id` - The ID of the subscription Alias.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the subscription Alias.
* `read` - (Defaults to 5 minutes) Used when retrieving the subscription Alias.
* `update` - (Defaults to 30 minutes) Used when updating the subscription Alias.
* `delete` - (Defaults to 30 minutes) Used when deleting the subscription Alias.

## Import

subscription Aliass can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_subscription_alias.example /providers/Microsoft.Subscription/aliases/alias1
```
