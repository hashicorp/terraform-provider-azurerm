---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_api_operation_tag"
description: |-
  Manages a API Management API Operation Tag.
---

# azurerm_api_management_api_operation_tag

Manages a API Management API Operation Tag.

## Example Usage

```hcl
data "azurerm_api_management_api" "example" {
  name                = "search-api"
  api_management_name = "search-api-management"
  resource_group_name = "search-service"
  revision            = "2"
}

resource "azurerm_api_management_api_operation" "example" {
  operation_id        = "user-delete"
  api_name            = data.azurerm_api_management_api.example.name
  api_management_name = data.azurerm_api_management_api.example.api_management_name
  resource_group_name = data.azurerm_api_management_api.example.resource_group_name
  display_name        = "Delete User Operation"
  method              = "DELETE"
  url_template        = "/users/{id}/delete"
  description         = "This can only be done by the logged in user."

  response {
    status_code = 200
  }
}

resource "azurerm_api_management_api_operation_tag" "example" {
  api_operation_id = azurerm_api_management_api_operation.example.id
  name             = "example-Tag"
}
```

## Arguments Reference

The following arguments are supported:

* `api_operation_id` - (Required) The ID of the API Management API Operation. Changing this forces a new API Management API Operation Tag to be created.

* `name` - (Required) The name which should be used for this API Management API Operation Tag. Changing this forces a new API Management API Operation Tag to be created. The name must be unique in the API Management Service.

---

* `display_name` - (Optional) The display name of the API Management API Operation Tag.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the API Management API Operation Tag.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Management API Operation Tag.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management API Operation Tag.
* `update` - (Defaults to 30 minutes) Used when updating the API Management API Operation Tag.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Management API Operation Tag.

## Import

API Management API Operation Tags can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_api_operation_tag.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ApiManagement/service/service1/apis/api1/operations/operation1/tags/tag1
```
