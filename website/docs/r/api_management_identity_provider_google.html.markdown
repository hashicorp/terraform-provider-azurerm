---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_identity_provider_google"
description: |-
  Manages an API Management Google Identity Provider.
---

# azurerm_api_management_identity_provider_google

Manages an API Management Google Identity Provider.

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

resource "azurerm_api_management_identity_provider_google" "example" {
  resource_group_name = azurerm_resource_group.example.name
  api_management_name = azurerm_api_management.example.name
  client_id           = "00000000.apps.googleusercontent.com"
  client_secret       = "00000000000000000000000000000000"
}
```

## Argument Reference

The following arguments are supported:

* `api_management_name` - (Required) The Name of the API Management Service where this Google Identity Provider should be created. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The Name of the Resource Group where the API Management Service exists. Changing this forces a new resource to be created.

* `client_id` - (Required) Client Id for Google Sign-in.

* `client_secret` - (Required) Client secret for Google Sign-in.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the API Management Google Identity Provider.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Management Google Identity Provider.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management Google Identity Provider.
* `update` - (Defaults to 30 minutes) Used when updating the API Management Google Identity Provider.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Management Google Identity Provider.

## Import

API Management Google Identity Provider can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_identity_provider_google.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.ApiManagement/service/instance1/identityProviders/google
```
