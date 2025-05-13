---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_api_operation"
description: |-
  Manages an API Operation within an API Management Service.
---

# azurerm_api_management_api_operation

Manages an API Operation within an API Management Service.

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

  template_parameter {
    name     = "id"
    type     = "number"
    required = true
  }

  response {
    status_code = 200
  }
}
```

## Argument Reference

The following arguments are supported:

* `operation_id` - (Required) A unique identifier for this API Operation. Changing this forces a new resource to be created.

* `api_name` - (Required) The name of the API within the API Management Service where this API Operation should be created. Changing this forces a new resource to be created.

* `api_management_name` - (Required) The Name of the API Management Service where the API exists. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The Name of the Resource Group in which the API Management Service exists. Changing this forces a new resource to be created.

* `display_name` - (Required) The Display Name for this API Management Operation.

* `method` - (Required) The HTTP Method used for this API Management Operation, like `GET`, `DELETE`, `PUT` or `POST` - but not limited to these values.

* `url_template` - (Required) The relative URL Template identifying the target resource for this operation, which may include parameters.

---

* `description` - (Optional) A description for this API Operation, which may include HTML formatting tags.

* `request` - (Optional) A `request` block as defined below.

* `response` - (Optional) One or more `response` blocks as defined below.

* `template_parameter` - (Optional) One or more `template_parameter` blocks as defined below. Required if `url_template` contains one or more parameters.

---

An `example` block supports the following:

* `name` - (Required) The name of this example.

* `summary` - (Optional) A short description for this example.

* `description` - (Optional) A long description for this example.

* `value` - (Optional) The example of the representation.

* `external_value` - (Optional) A URL that points to the literal example.

---

A `form_parameter` block supports the following:

* `name` - (Required) The Name of this Form Parameter.

* `required` - (Required) Is this Form Parameter Required?

* `type` - (Required) The Type of this Form Parameter, such as a `string`.

* `description` - (Optional) A description of this Form Parameter.

* `default_value` - (Optional) The default value for this Form Parameter.

* `values` - (Optional) One or more acceptable values for this Form Parameter.

* `example` - (Optional) One or more `example` blocks as defined above.

* `schema_id` - (Optional) The name of the Schema.

* `type_name` - (Optional) The type name defined by the Schema.

---

A `header` block supports the following:

* `name` - (Required) The Name of this Header.

* `required` - (Required) Is this Header Required?

* `type` - (Required) The Type of this Header, such as a `string`.

* `description` - (Optional) A description of this Header.

* `default_value` - (Optional) The default value for this Header.

* `values` - (Optional) One or more acceptable values for this Header.

* `example` - (Optional) One or more `example` blocks as defined above.

* `schema_id` - (Optional) The name of the Schema.

* `type_name` - (Optional) The type name defined by the Schema.

---

A `query_parameter` block supports the following:

* `name` - (Required) The Name of this Query Parameter.

* `required` - (Required) Is this Query Parameter Required?

* `type` - (Required) The Type of this Query Parameter, such as a `string`.

* `description` - (Optional) A description of this Query Parameter.

* `default_value` - (Optional) The default value for this Query Parameter.

* `values` - (Optional) One or more acceptable values for this Query Parameter.

* `example` - (Optional) One or more `example` blocks as defined above.

* `schema_id` - (Optional) The name of the Schema.

* `type_name` - (Optional) The type name defined by the Schema.

---

A `request` block supports the following:

* `description` - (Optional) A description of the HTTP Request, which may include HTML tags.

* `header` - (Optional) One or more `header` blocks as defined above.

* `query_parameter` - (Optional) One or more `query_parameter` blocks as defined above.

* `representation` - (Optional) One or more `representation` blocks as defined below.

---

A `representation` block supports the following:

* `content_type` - (Required) The Content Type of this representation, such as `application/json`.

* `form_parameter` - (Optional) One or more `form_parameter` block as defined above.

-> **Note:** This is Required when `content_type` is set to `application/x-www-form-urlencoded` or `multipart/form-data`.

* `example` - (Optional) One or more `example` blocks as defined above.

* `schema_id` - (Optional) The ID of an API Management Schema which represents this Response.

-> **Note:** This can only be specified when `content_type` is not set to `application/x-www-form-urlencoded` or `multipart/form-data`.

* `type_name` - (Optional) The Type Name defined by the Schema.

-> **Note:** This can only be specified when `content_type` is not set to `application/x-www-form-urlencoded` or `multipart/form-data`.

---

A `response` block supports the following:

* `status_code` - (Required) The HTTP Status Code.

* `description` - (Optional) A description of the HTTP Response, which may include HTML tags.

* `header` - (Optional) One or more `header` blocks as defined above.

* `representation` - (Optional) One or more `representation` blocks as defined above.

---

A `template_parameter` block supports the following:

* `name` - (Required) The Name of this Template Parameter.

* `required` - (Required) Is this Template Parameter Required?

* `type` - (Required) The Type of this Template Parameter, such as a `string`.

* `description` - (Optional) A description of this Template Parameter.

* `default_value` - (Optional) The default value for this Template Parameter.

* `values` - (Optional) One or more acceptable values for this Template Parameter.

* `example` - (Optional) One or more `example` blocks as defined above.

* `schema_id` - (Optional) The name of the Schema.

* `type_name` - (Optional) The type name defined by the Schema.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the API Management API Operation.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Management API Operation.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management API Operation.
* `update` - (Defaults to 30 minutes) Used when updating the API Management API Operation.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Management API Operation.

## Import

API Management API Operation's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_api_operation.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.ApiManagement/service/instance1/apis/api1/operations/operation1
```
