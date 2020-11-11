---
subcategory: "Base"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_subscription_alias"
description: |-
  Gets information about an existing Subscription Alias.
---

# Data Source: azurerm_subscription_alias

Use this data source to access information about an existing Subscription Alias.

## Example Usage

```hcl
data "azurerm_subscription_alias" "example" {
  name = "example-alias"
}

output "id" {
  value = data.azurerm_subscription_alias.example.subscription_id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this subscription Alias.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the subscription Alias.

* `subscription_id` - The subscription ID which this Subscription Alias is associated to.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the subscription Alias.
