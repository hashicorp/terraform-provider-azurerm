---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_product_tag"
description: |-
  Manages an API Management Product tag
---

# azurerm_api_management_product_tag

Manages an API Management Product tag

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

resource "azurerm_api_management_product" "example" {
  product_id            = "test-product"
  api_management_name   = azurerm_api_management.example.name
  resource_group_name   = azurerm_resource_group.example.name
  display_name          = "Test Product"
  subscription_required = true
  approval_required     = true
  published             = true
}

resource "azurerm_api_management_tag" "example" {
  api_management_id = azurerm_api_management.example.id
  name              = "example-tag"
}

resource "azurerm_api_management_product_tag" "example" {
  api_management_product_id = azurerm_api_management_product.example.product_id
  api_management_name       = azurerm_api_management.example.name
  resource_group_name       = azurerm_resource_group.example.name
  name                      = azurerm_api_management_tag.example.name
}
```

## Argument Reference

The following arguments are supported:

* `api_management_name` - (Required) The name of the API Management Service. Changing this forces a new resource to be created.

* `api_management_product_id` - (Required) The name of the API Management product. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which the API Management Service should be exist. Changing this forces a new resource to be created.

* `name` - (Required) The name which should be used for this API Management Tag. Changing this forces a new API Management Tag to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the API Management Product.

## Timeouts

The `timeouts` block allows you to
specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Management Product.
* `update` - (Defaults to 30 minutes) Used when updating the API Management Product.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management Product.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Management Product.

## Import

API Management Products can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_product_tag.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.ApiManagement/service/instance1/products/myproduct/tags/mytag
```
