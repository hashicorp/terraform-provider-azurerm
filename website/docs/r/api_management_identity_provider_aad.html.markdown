---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_identity_provider_aad"
sidebar_current: "docs-azurerm-resource-api-management-identity-provider-aad"
description: |-
  Manages an API Management AAD Identity Provider.
---

# azurerm_api_management_identity_provider_aad

Manages an API Management AAD Identity Provider.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_api_management" "test" {
  name                = "example-apim"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  publisher_name      = "My Company"
  publisher_email     = "company@terraform.io"
  sku_name            = "Developer_1"
}

resource "azurerm_api_management_identity_provider_aad" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  api_management_name = "${azurerm_api_management.test.name}"
  client_id           = "aadclientid"
  client_secret       = "aadsecret"
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

---

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the API Management AAD Identity Provider.

## Import

API Management AAD Identity Provider can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_identity_provider_aad.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.ApiManagement/service/instance1/identityProviders/aad
```