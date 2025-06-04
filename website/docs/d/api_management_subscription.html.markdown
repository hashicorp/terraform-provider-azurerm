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
  api_management_id = "example-apim"
  subscription_id   = "example-subscription-id"
}

output "id" {
  value = data.azurerm_api_management_subscription.example.subscription_id
}
```

## Arguments Reference

The following arguments are supported:

* `api_management_id` - (Required) The ID of the API Management Service in which this Subscription exists.

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

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.ApiManagement`: 2022-08-01
