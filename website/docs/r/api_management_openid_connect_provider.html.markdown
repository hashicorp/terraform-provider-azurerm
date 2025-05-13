---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_openid_connect_provider"
description: |-
  Manages an OpenID Connect Provider within a API Management Service.
---

# azurerm_api_management_openid_connect_provider

Manages an OpenID Connect Provider within a API Management Service.

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

  sku_name = "Developer_1"
}

resource "azurerm_api_management_openid_connect_provider" "example" {
  name                = "example-provider"
  api_management_name = azurerm_api_management.example.name
  resource_group_name = azurerm_resource_group.example.name
  client_id           = "00001111-2222-3333-4444-555566667777"
  client_secret       = "00001111-423egvwdcsjx-00001111"
  display_name        = "Example Provider"
  metadata_endpoint   = "https://example.com/example"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) the Name of the OpenID Connect Provider which should be created within the API Management Service. Changing this forces a new resource to be created.

* `api_management_name` - (Required) The name of the API Management Service in which this OpenID Connect Provider should be created. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the API Management Service exists. Changing this forces a new resource to be created.

* `client_id` - (Required) The Client ID used for the Client Application.

* `client_secret` - (Required) The Client Secret used for the Client Application.

* `display_name` - (Required) A user-friendly name for this OpenID Connect Provider.

* `metadata_endpoint` - (Required) The URI of the Metadata endpoint.

---

* `description` - (Optional) A description of this OpenID Connect Provider.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the API Management OpenID Connect Provider.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Management OpenID Connect Provider.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management OpenID Connect Provider.
* `update` - (Defaults to 30 minutes) Used when updating the API Management OpenID Connect Provider.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Management OpenID Connect Provider.

## Import

API Management OpenID Connect Providers can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_openid_connect_provider.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.ApiManagement/service/instance1/openidConnectProviders/provider1
```
