---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_identity_provider_aadb2c"
description: |-
  Manages an API Management Azure AD B2C Identity Provider.
---

# azurerm_api_management_identity_provider_aadb2c

Manages an API Management Azure AD B2C Identity Provider.

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

resource "azuread_application" "example" {
  display_name = "acctestam-example"
}

resource "azuread_application_password" "example" {
  application_object_id = azuread_application.example.object_id
  end_date_relative     = "36h"
}

resource "azurerm_api_management_identity_provider_aadb2c" "example" {
  resource_group_name = azurerm_resource_group.example.name
  api_management_name = azurerm_api_management.example.name
  client_id           = azuread_application.example.application_id
  client_secret       = "P@55w0rD!"
  allowed_tenant      = "myb2ctenant.onmicrosoft.com"
  signin_tenant       = "myb2ctenant.onmicrosoft.com"
  authority           = "myb2ctenant.b2clogin.com"
  signin_policy       = "B2C_1_Login"
  signup_policy       = "B2C_1_Signup"

  depends_on = [azuread_application_password.example]
}
```

## Arguments Reference

The following arguments are supported:

* `api_management_name` - (Required) The Name of the API Management Service where this AAD Identity Provider should be created. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The Name of the Resource Group where the API Management Service exists. Changing this forces a new resource to be created.

* `client_id` - (Required) Client ID of the Application in your B2C tenant.

* `client_secret` - (Required) Client secret of the Application in your B2C tenant.

* `allowed_tenant` - (Required) The allowed AAD tenant, usually your B2C tenant domain.

* `signin_tenant` - (Required) The tenant to use instead of Common when logging into Active Directory, usually your B2C tenant domain.

* `authority` - (Required) OpenID Connect discovery endpoint hostname, usually your b2clogin.com domain.

* `signin_policy` - (Required) Signin Policy Name.

* `signup_policy` - (Required) Signup Policy Name.

* `client_library` - (Optional) The client library to be used in the Azure AD B2C Identity Provider.

---

* `password_reset_policy` - (Optional) Password reset Policy Name.

* `profile_editing_policy` - (Optional) Profile editing Policy Name.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the API Management Azure AD B2C Identity Provider Resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Management Azure AD B2C Identity Provider.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management Azure AD B2C Identity Provider.
* `update` - (Defaults to 30 minutes) Used when updating the API Management Azure AD B2C Identity Provider.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Management Azure AD B2C Identity Provider.

## Import

API Management Azure AD B2C Identity Providers can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_identity_provider_aadb2c.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.ApiManagement/service/service1/identityProviders/aadB2C
```
