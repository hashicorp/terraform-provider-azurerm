---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_product"
description: |-
  Gets information about an existing API Management Product.
---

# Data Source: azurerm_api_management_product

Use this data source to access information about an existing API Management Product.

## Example Usage

```hcl
data "azurerm_api_management_product" "example" {
  product_id          = "my-product"
  api_management_name = "example-apim"
  resource_group_name = "search-service"
}

output "product_terms" {
  value = data.azurerm_api_management_product.example.terms
}
```

## Argument Reference

* `api_management_name` - The Name of the API Management Service in which this Product exists.

* `product_id` - The Identifier for the API Management Product.

* `resource_group_name` - The Name of the Resource Group in which the API Management Service exists.

## Attributes Reference

* `id` - The ID of the API Management Product.

* `approval_required` - Do subscribers need to be approved prior to being able to use the Product?

* `display_name` - The Display Name for this API Management Product.

* `published` - Is this Product Published?

* `subscription_required` - Is a Subscription required to access API's included in this Product?

* `description` - The description of this Product, which may include HTML formatting tags.

* `subscriptions_limit` - The number of subscriptions a user can have to this Product at the same time.

* `terms` - Any Terms and Conditions for this Product, which must be accepted by Developers before they can begin the Subscription process.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the API Management Product.
