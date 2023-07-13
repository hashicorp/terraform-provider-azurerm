---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_product_group"
description: |-
  Manages an API Management Product Assignment to a Group.
---

# azurerm_api_management_product_group

Manages an API Management Product Assignment to a Group.

## Example Usage

```hcl
data "azurerm_api_management" "example" {
  name                = "example-api"
  resource_group_name = "example-resources"
}

data "azurerm_api_management_product" "example" {
  product_id          = "my-product"
  api_management_name = data.azurerm_api_management.example.name
  resource_group_name = data.azurerm_api_management.example.resource_group_name
}

data "azurerm_api_management_group" "example" {
  name                = "my-group"
  api_management_name = data.azurerm_api_management.example.name
  resource_group_name = data.azurerm_api_management.example.resource_group_name
}

resource "azurerm_api_management_product_group" "example" {
  product_id          = data.azurerm_api_management_product.example.product_id
  group_name          = data.azurerm_api_management_group.example.name
  api_management_name = data.azurerm_api_management.example.name
  resource_group_name = data.azurerm_api_management.example.resource_group_name
}
```

## Argument Reference

The following arguments are supported:

* `product_id` - (Required) The ID of the API Management Product within the API Management Service. Changing this forces a new resource to be created.

* `group_name` - (Required) The Name of the API Management Group within the API Management Service. Changing this forces a new resource to be created.

* `api_management_name` - (Required) The name of the API Management Service. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which the API Management Service exists. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the API Management Product Group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Management Product Group.
* `update` - (Defaults to 30 minutes) Used when updating the API Management Product Group.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management Product Group.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Management Product Group.

## Import

API Management Product Groups can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_product_group.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ApiManagement/service/service1/products/exampleId/groups/groupId
```
