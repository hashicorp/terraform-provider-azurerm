---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_api_management_subscription"
description: |-
  Gets information about an existing API Management Subscription.
---

# Data Source: azurerm_api_management_subscription

Use this data source to access information about an existing API Management Subscription.

## Example Usage

```hcl
data "azurerm_api_management_subscription" "example" {
  resource_group_name = "example-rg"
  subscription_id     = "example-subscription"
  api_management_name = "example-apim"
}

output "id" {
  value = data.azurerm_api_management_subscription.example.subscription_id
}
```

## Arguments Reference

The following arguments are supported:

* `api_management_name` - (Required) The Name of the API Management Service in which this Subscription exists.

* `resource_group_name` - (Required) The Name of the Resource Group in which the API Management Service exists.

* `subscription_id` - (Required) The Identifier for the API Management Subscription.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the API Management Subscription.

* `allow_tracing` - Indicates whether tracing is enabled.

* `api_id` - The ID of the API assigned to this Subscription.

* `display_name` - The display name of this Subscription.

* `primary_key` - The primary key for this subscription.

* `product_id` - The ID of the Product assigned to this Subscription.

* `secondary_key` - The secondary key for this subscription.

* `state` - The state of this Subscription.

* `user_id` - The ID of the User assigned to this Subscription.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the API Management Subscription.