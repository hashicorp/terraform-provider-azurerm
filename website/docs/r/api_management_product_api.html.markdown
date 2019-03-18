---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_product_api"
sidebar_current: "docs-azurerm-resource-api-management-product-api"
description: |-
  Manages an API Management API Assignment to a Product.
---

# azurerm_api_management_product_api

Manages an API Management API Assignment to a Product.

## Example Usage

```hcl
data "azurerm_api_management" "example" {
  name                = "example-api"
  resource_group_name = "example-resources"
}


data "azurerm_api_management_api" "example" {
  name                = "search-api"
  api_management_name = "${data.azurerm_api_management.example.name}"
  resource_group_name = "${data.azurerm_api_management.example.resource_group_name}"
  revision            = "2"
}
data "azurerm_api_management_product" "test" {
  product_id          = "my-product"
  api_management_name = "${data.azurerm_api_management.example.name}"
  resource_group_name = "${data.azurerm_api_management.example.resource_group_name}"
}

resource "azurerm_api_management_product_api" "example" {
  api_name            = "${data.azurerm_api_management_api.example.name}"
  product_id          = "${data.azurerm_api_management_product.example.product_id}"
  api_management_name = "${data.azurerm_api_management.example.name}"
  resource_group_name = "${data.azurerm_api_management.example.resource_group_name}"
}
```

## Argument Reference

The following arguments are supported:

* `api_name` - (Required) The Name of the API Management API within the API Management Service. Changing this forces a new resource to be created.

* `api_management_name` - (Required) The name of the API Management Service. Changing this forces a new resource to be created.

* `product_id` - (Required) The ID of the API Management Product within the API Management Service. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which the API Management Service exists. Changing this forces a new resource to be created.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the API Management Product API.

## Import

API Management Product API's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_product_api.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ApiManagement/service/service1/products/exampleId/apis/apiId
```
