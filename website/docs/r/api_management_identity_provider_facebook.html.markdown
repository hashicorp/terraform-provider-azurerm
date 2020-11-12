---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_identity_provider_facebook"
description: |-
  Manages an API Management Facebook Identity Provider.
---

# azurerm_api_management_identity_provider_facebook

Manages an API Management Facebook Identity Provider.

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

resource "azurerm_api_management_identity_provider_facebook" "example" {
  api_management_id = azurerm_api_management.example.id
  app_id            = "00000000000000000000000000000000"
  app_secret        = "00000000000000000000000000000000"
}
```

## Argument Reference

The following arguments are supported:

* `api_management_id` - (Required) The ID of the API Management Service where this Facebook Identity Provider should be created. Changing this forces a new resource to be created.

* `api_management_name` - (Optional, Deprecated) The Name of the API Management Service where this Facebook Identity Provider should be created. This property is deprecated and will be removed in version 3.0 of the provider. Use the `api_management_id` property instead. Changing this forces a new resource to be created.

* `resource_group_name` - (Optional, Deprecated) The Name of the Resource Group where the API Management Service exists. This property is deprecated and will be removed in version 3.0 of the provider. Use the `api_management_id` property instead. Changing this forces a new resource to be created.

* `app_id` - (Required) App ID for Facebook.

* `app_secret` - (Required) App Secret for Facebook.

---

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the API Management Facebook Identity Provider.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Management Facebook Identity Provider.
* `update` - (Defaults to 30 minutes) Used when updating the API Management Facebook Identity Provider.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management Facebook Identity Provider.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Management Facebook Identity Provider.

## Import

API Management Facebook Identity Provider can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_identity_provider_facebook.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.ApiManagement/service1/identityProviders/Facebook
```
