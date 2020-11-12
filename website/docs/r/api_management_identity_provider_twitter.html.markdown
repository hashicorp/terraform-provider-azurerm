---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_identity_provider_twitter"
description: |-
  Manages an API Management Twitter Identity Provider.
---

# azurerm_api_management_identity_provider_twitter

Manages an API Management Twitter Identity Provider.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_api_management" "example" {
  name                = "example-apim"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  publisher_name      = "My Company"
  publisher_email     = "company@terraform.io"
  sku_name            = "Developer_1"
}

resource "azurerm_api_management_identity_provider_twitter" "example" {
  api_management_id = azurerm_api_management.example.id
  api_key           = "00000000000000000000000000000000"
  api_secret_key    = "00000000000000000000000000000000"
}
```

## Argument Reference

The following arguments are supported:

* `api_management_id` - (Required) The ID of the API Management Service where this Twitter Identity Provider should be created. Changing this forces a new resource to be created.

* `api_management_name` - (Optional, Deprecated) The Name of the API Management Service where this Twitter Identity Provider should be created. This property is deprecated and will be removed in version 3.0 of the provider. Use the `api_management_id` property instead. Changing this forces a new resource to be created.

* `resource_group_name` - (Optional, Deprecated) The Name of the Resource Group where the API Management Service exists. This property is deprecated and will be removed in version 3.0 of the provider. Use the `api_management_id` property instead. Changing this forces a new resource to be created.

* `api_key` - (Required) App Consumer API key for Twitter.

* `api_secret_key` - (Required) App Consumer API secret key for Twitter.

---

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the API Management Twitter Identity Provider.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Management Twitter Identity Provider.
* `update` - (Defaults to 30 minutes) Used when updating the API Management Twitter Identity Provider.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management Twitter Identity Provider.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Management Twitter Identity Provider.

## Import

API Management Twitter Identity Provider can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_identity_provider_twitter.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.ApiManagement/service1/identityProviders/Twitter
```
