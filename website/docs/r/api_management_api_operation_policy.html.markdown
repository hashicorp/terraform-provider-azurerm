---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_api_operation_policy"
description: |-
  Manages an API Management API Operation Policy
---

# azurerm_api_management_api_operation_policy

Manages an API Management API Operation Policy

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

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

resource "azurerm_api_management_api" "example" {
  name                = "example-api"
  resource_group_name = azurerm_resource_group.example.name
  api_management_name = azurerm_api_management.example.name
  revision            = "1"
}

resource "azurerm_api_management_api_operation" "example" {
  operation_id        = "acctest-operation"
  api_name            = azurerm_api_management_api.example.name
  api_management_name = azurerm_api_management.example.name
  resource_group_name = azurerm_resource_group.example.name
  display_name        = "DELETE Resource"
  method              = "DELETE"
  url_template        = "/resource"
}

resource "azurerm_api_management_api_operation_policy" "example" {
  api_name            = azurerm_api_management_api_operation.example.api_name
  api_management_name = azurerm_api_management_api_operation.example.api_management_name
  resource_group_name = azurerm_api_management_api_operation.example.resource_group_name
  operation_id        = azurerm_api_management_api_operation.example.operation_id

  xml_content = <<XML
<policies>
  <inbound>
    <find-and-replace from="xyz" to="abc" />
  </inbound>
</policies>
XML

}
```

## Argument Reference

The following arguments are supported:

* `api_name` - (Required) The name of the API within the API Management Service where the Operation exists. Changing this forces a new resource to be created.

* `api_management_name` - (Required) The name of the API Management Service. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which the API Management Service exists. Changing this forces a new resource to be created.

* `operation_id` - (Required) The operation identifier within an API. Must be unique in the current API Management service instance. Changing this forces a new resource to be created.

* `xml_content` - (Optional) The XML Content for this Policy.

* `xml_link` - (Optional) A link to a Policy XML Document, which must be publicly available.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the API Management API Operation Policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Management API Operation Policy.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management API Operation Policy.
* `update` - (Defaults to 30 minutes) Used when updating the API Management API Operation Policy.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Management API Operation Policy.

## Import

API Management API Operation Policy can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_api_operation_policy.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ApiManagement/service/instance1/apis/api1/operations/operation1
```
