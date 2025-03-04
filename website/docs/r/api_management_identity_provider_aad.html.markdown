---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_identity_provider_aad"
description: |-
  Manages an API Management AAD Identity Provider.
---

# azurerm_api_management_identity_provider_aad

Manages an API Management AAD Identity Provider.

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

resource "azurerm_api_management_identity_provider_aad" "example" {
  resource_group_name = azurerm_resource_group.example.name
  api_management_name = azurerm_api_management.example.name
  client_id           = "00000000-0000-0000-0000-000000000000"
  client_secret       = "00000000000000000000000000000000"
  allowed_tenants     = ["00000000-0000-0000-0000-000000000000"]
}
```

## Argument Reference

The following arguments are supported:

* `api_management_name` - (Required) The Name of the API Management Service where this AAD Identity Provider should be created. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The Name of the Resource Group where the API Management Service exists. Changing this forces a new resource to be created.

* `client_id` - (Required) Client Id of the Application in the AAD Identity Provider.

* `client_secret` - (Required) Client secret of the Application in the AAD Identity Provider.

* `allowed_tenants` - (Required) List of allowed AAD Tenants.

* `client_library` - (Optional) The client library to be used in the AAD Identity Provider.

* `signin_tenant` - (Optional) The AAD Tenant to use instead of Common when logging into Active Directory.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the API Management AAD Identity Provider.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Management AAD Identity Provider.
* `update` - (Defaults to 30 minutes) Used when updating the API Management AAD Identity Provider.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management AAD Identity Provider.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Management AAD Identity Provider.

## Import

API Management AAD Identity Provider can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_identity_provider_aad.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.ApiManagement/service/instance1/identityProviders/aad
```
