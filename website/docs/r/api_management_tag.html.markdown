---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_tag"
description: |-
  Manages a API Management Tag.
---

# azurerm_api_management_tag

Manages a API Management Tag.

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
  sku_name            = "Consumption_0"
}

resource "azurerm_api_management_tag" "example" {
  api_management_id = azurerm_api_management.example.id
  name              = "example-Tag"
}
```

## Arguments Reference

The following arguments are supported:

* `api_management_id` - (Required) The ID of the API Management. Changing this forces a new API Management Tag to be created.

* `name` - (Required) The name which should be used for this API Management Tag. Changing this forces a new API Management Tag to be created. The name must be unique in the API Management Service.

* `display_name` - (Optional) The display name of the API Management Tag. Defaults to the `name`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the API Management Tag.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Management Tag.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management Tag.
* `update` - (Defaults to 30 minutes) Used when updating the API Management Tag.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Management Tag.

## Import

API Management Tags can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_tag.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ApiManagement/service/service1/tags/tag1
```
