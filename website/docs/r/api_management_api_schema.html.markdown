---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_api_schema"
description: |-
  Manages an API Schema within an API Management Service.
---

# azurerm_api_management_api_schema

Manages an API Schema within an API Management Service.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

data "azurerm_api_management_api" "example" {
  name                = "search-api"
  api_management_name = "search-api-management"
  resource_group_name = "search-service"
  revision            = "2"
}

resource "azurerm_api_management_api_schema" "example" {
  api_name            = data.azurerm_api_management_api.example.name
  api_management_name = data.azurerm_api_management_api.example.api_management_name
  resource_group_name = data.azurerm_api_management_api.example.resource_group_name
  schema_id           = "example-schema"
  content_type        = "application/vnd.ms-azure-apim.xsd+xml"
  value               = file("api_management_api_schema.xml")
}
```

## Argument Reference

The following arguments are supported:

* `schema_id` - (Required) A unique identifier for this API Schema. Changing this forces a new resource to be created.

* `api_name` - (Required) The name of the API within the API Management Service where this API Schema should be created. Changing this forces a new resource to be created.

* `api_management_name` - (Required) The Name of the API Management Service where the API exists. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The Name of the Resource Group in which the API Management Service exists. Changing this forces a new resource to be created.

* `content_type` - (Required) The content type of the API Schema.

* `value` - (Optional) The JSON escaped string defining the document representing the Schema.

* `components` - (Optional) Types definitions. Used for Swagger/OpenAPI v2/v3 schemas only.

* `definitions` - (Optional) Types definitions. Used for Swagger/OpenAPI v1 schemas only.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the API Management API Schema.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Management API Schema.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management API Schema.
* `update` - (Defaults to 30 minutes) Used when updating the API Management API Schema.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Management API Schema.

## Import

API Management API Schema's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_api_schema.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.ApiManagement/service/instance1/apis/api1/schemas/schema1
```
