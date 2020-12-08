---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_subscription"
description: |-
  Manages a Subscription within a API Management Service.
---

# azurerm_api_management_subscription

Manages a Subscription within a API Management Service.


## Example Usage

```hcl
data "azurerm_api_management" "example" {
  name                = "example-apim"
  resource_group_name = "example-resources"
}

data "azurerm_api_management_product" "example" {
  product_id          = "00000000-0000-0000-0000-000000000000"
  api_management_name = data.azurerm_api_management.example.name
  resource_group_name = data.azurerm_api_management.example.resource_group_name
}

data "azurerm_api_management_user" "example" {
  user_id             = "11111111-1111-1111-1111-111111111111"
  api_management_name = data.azurerm_api_management.example.name
  resource_group_name = data.azurerm_api_management.example.resource_group_name
}

resource "azurerm_api_management_subscription" "example" {
  api_management_name = data.azurerm_api_management.example.name
  resource_group_name = data.azurerm_api_management.example.resource_group_name
  user_id             = data.azurerm_api_management_user.example.id
  product_id          = data.azurerm_api_management_product.example.id
  display_name        = "Parser API"
}
```


## Argument Reference

The following arguments are supported:

* `api_management_name` - (Required) The name of the API Management Service where this Subscription should be created. Changing this forces a new resource to be created.

* `display_name` - (Required) The display name of this Subscription.

* `resource_group_name` - (Required) The name of the Resource Group in which the API Management Service exists. Changing this forces a new resource to be created.

* `product_id` - (Required) The ID of the Product which should be assigned to this Subscription. Changing this forces a new resource to be created.

* `user_id` - (Required) The ID of the User which should be assigned to this Subscription. Changing this forces a new resource to be created.

---

* `state` - (Optional) The state of this Subscription. Possible values are `active`, `cancelled`, `expired`, `rejected`, `submitted` and `suspended`. Defaults to `submitted`.

* `subscription_id` - (Optional) An Identifier which should used as the ID of this Subscription. If not specified a new Subscription ID will be generated. Changing this forces a new resource to be created.

* `allow_tracing` - (Optional) Determines whether tracing can be enabled.  Defaults to `true`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the API Management Subscription.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Management Subscription.
* `update` - (Defaults to 30 minutes) Used when updating the API Management Subscription.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management Subscription.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Management Subscription.

## Import

API Management Subscriptions can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_subscription.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-resources/providers/Microsoft.ApiManagement/service/example-apim/subscriptions/subscription-name
```
