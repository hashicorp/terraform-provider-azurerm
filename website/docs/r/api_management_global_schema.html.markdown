---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_global_schema"
description: |-
  Manages a Global Schema within an API Management Service.
---

# azurerm_api_management_global_schema

Manages a Global Schema within an API Management Service.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "West Europe"
}

resource "azurerm_api_management" "example" {
  name                = "example-apim"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Consumption_0"
}

resource "azurerm_api_management_global_schema" "example" {
  schema_id           = "example-schema1"
  api_management_name = azurerm_api_management.example.name
  resource_group_name = azurerm_resource_group.example.name
  type                = "xml"
  value               = file("api_management_api_schema.xml")
}
```

## Argument Reference

The following arguments are supported:

* `schema_id` - (Required) A unique identifier for this Schema. Changing this forces a new resource to be created.

* `api_management_name` - (Required) The Name of the API Management Service where the API exists. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The Name of the Resource Group in which the API Management Service exists. Changing this forces a new resource to be created.

* `type` - (Required) The content type of the Schema. Possible values are `xml` and `json`.

* `value` - (Required) The string defining the document representing the Schema.

* `description` - (Optional) The description of the schema.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the API Management API Schema.

## Timeouts

The `timeouts` block allows you to
specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Management API Schema.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management API Schema.
* `update` - (Defaults to 30 minutes) Used when updating the API Management API Schema.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Management API Schema.

## Import

API Management API Schema's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_global_schema.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.ApiManagement/service/instance1/schemas/schema1
```
