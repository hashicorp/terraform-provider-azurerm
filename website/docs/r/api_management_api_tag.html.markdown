---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_api_tag"
description: |-
  Manages an API Management API Tag.
---

# azurerm_api_management_api_tag

Manages the Assignment of an API Management API Tag to an API.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

data "azurerm_api_management" "example" {
  name                = "example-apim"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_api_management_api" "example" {
  name                = "example-api"
  resource_group_name = azurerm_resource_group.example.name
  api_management_name = data.azurerm_api_management.example.name
  revision            = "1"
}

resource "azurerm_api_management_tag" "example" {
  api_management_id = data.azurerm_api_management.example.id
  name              = "example-tag"
}

resource "azurerm_api_management_api_tag" "example" {
  api_id = azurerm_api_management_api.example.id
  name   = "example-tag"
}
```

## Arguments Reference

The following arguments are supported:

* `api_id` - (Required) The ID of the API Management API. Changing this forces a new API Management API Tag to be created.

* `name` - (Required) The name of the tag. It must be known in the API Management instance. Changing this forces a new API Management API Tag to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the API Management API Tag.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Management API Tag.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management API Tag.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Management API Tag.

## Import

API Management API Tags can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_api_tag.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ApiManagement/service/service1/apis/api1/tags/tag1
```
