---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_identity_provider_aadb2c"
description: |-
  Manages a API Management Resources.
---

# azurerm_api_management_identity_provider_aadb2c

Manages a API Management Resources.

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

resource "azurerm_api_management_identity_provider_aadb2c" "example" {
  resource_group_name = azurerm_resource_group.example.name
  api_management_name = azurerm_api_management.example.name
  client_id           = "00000000-0000-0000-0000-000000000000"
  client_secret       = "00000000000000000000000000000000"
  signin_tenant       = "00000000-0000-0000-0000-000000000000"
  authority           = "ExampleAuthority"
  signin_policy       = "ExampleSigninPolicy"
  signup_policy       = "ExampleSignupPolicy"
}
```

## Arguments Reference

The following arguments are supported:

* `api_management_name` - (Required) The Name of the API Management Service where this AADB2C Identity Provider should be created. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The Name of the Resource Group where the API Management Service exists. Changing this forces a new resource to be created.

* `client_id` - (Required) Client Id of the Application in the AADB2C Identity Provider.

* `client_secret` - (Required) Client secret of the Application in the AADB2C Identity Provider.

* `allowed_tenants` - (Required) List of allowed AAD Tenants.

* `authority` - (Required) TODO.

* `signin_policy` - (Required) TODO.

* `signin_tenant` - (Required) TODO.

* `signup_policy` - (Required) TODO.

---

* `password_reset_policy` - (Optional) TODO.

* `profile_editing_policy` - (Optional) TODO.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the API Management Resources.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Management Resources.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management Resources.
* `update` - (Defaults to 30 minutes) Used when updating the API Management Resources.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Management Resources.

## Import

API Management Resourcess can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_identity_provider_aadb2c.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.ApiManagement/service/instance1/identityProviders/aadb2c
```
