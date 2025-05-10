---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_authorization_server"
description: |-
  Manages an Authorization Server within an API Management Service.
---

# azurerm_api_management_authorization_server

Manages an Authorization Server within an API Management Service.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

data "azurerm_api_management" "example" {
  name                = "search-api"
  resource_group_name = "search-service"
}

resource "azurerm_api_management_authorization_server" "example" {
  name                         = "test-server"
  api_management_name          = data.azurerm_api_management.example.name
  resource_group_name          = data.azurerm_api_management.example.resource_group_name
  display_name                 = "Test Server"
  authorization_endpoint       = "https://example.mydomain.com/client/authorize"
  client_id                    = "42424242-4242-4242-4242-424242424242"
  client_registration_endpoint = "https://example.mydomain.com/client/register"

  grant_types = [
    "authorizationCode",
  ]
  authorization_methods = [
    "GET",
  ]
}
```

## Argument Reference

The following arguments are supported:

* `api_management_name` - (Required) The name of the API Management Service in which this Authorization Server should be created. Changing this forces a new resource to be created.

* `authorization_methods` - (Required) The HTTP Verbs supported by the Authorization Endpoint. Possible values are `DELETE`, `GET`, `HEAD`, `OPTIONS`, `PATCH`, `POST`, `PUT` and `TRACE`.

-> **Note:** `GET` must always be present.

* `authorization_endpoint` - (Required) The OAUTH Authorization Endpoint.

* `client_id` - (Required) The Client/App ID registered with this Authorization Server.

* `client_registration_endpoint` - (Required) The URI of page where Client/App Registration is performed for this Authorization Server.

* `display_name` - (Required) The user-friendly name of this Authorization Server.

* `grant_types` - (Required) Form of Authorization Grants required when requesting an Access Token. Possible values are `authorizationCode`, `clientCredentials`, `implicit` and `resourceOwnerPassword`.

* `name` - (Required) The name of this Authorization Server. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which the API Management Service exists. Changing this forces a new resource to be created.

---

* `bearer_token_sending_methods` - (Optional) The mechanism by which Access Tokens are passed to the API. Possible values are `authorizationHeader` and `query`.

* `client_authentication_method` - (Optional) The Authentication Methods supported by the Token endpoint of this Authorization Server.. Possible values are `Basic` and `Body`.

* `client_secret` - (Optional) The Client/App Secret registered with this Authorization Server.

* `default_scope` - (Optional) The Default Scope used when requesting an Access Token, specified as a string containing space-delimited values.

* `description` - (Optional) A description of the Authorization Server, which may contain HTML formatting tags.

* `resource_owner_password` - (Optional) The password associated with the Resource Owner.

-> **Note:** This can only be specified when `grant_type` includes `resourceOwnerPassword`.

* `resource_owner_username` - (Optional) The username associated with the Resource Owner.

-> **Note:** This can only be specified when `grant_type` includes `resourceOwnerPassword`.

* `support_state` - (Optional) Does this Authorization Server support State? If this is set to `true` the client may use the state parameter to raise protocol security.

* `token_body_parameter` - (Optional) A `token_body_parameter` block as defined below.

* `token_endpoint` - (Optional) The OAUTH Token Endpoint.

---

A `token_body_parameter` block supports the following:

* `name` - (Required) The Name of the Parameter.

* `value` - (Required) The Value of the Parameter.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the API Management Authorization Server.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Management Authorization Server.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management Authorization Server.
* `update` - (Defaults to 30 minutes) Used when updating the API Management Authorization Server.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Management Authorization Server.

## Import

API Management Authorization Servers can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_authorization_server.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ApiManagement/service/service1/authorizationServers/server1
```
