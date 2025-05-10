---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_gateway_api"
description: |-
  Manages a API Management Gateway API.
---

# azurerm_api_management_gateway_api

Manages a API Management Gateway API.

## Example Usage

```hcl
data "azurerm_api_management" "example" {
  name                = "example-api"
  resource_group_name = "example-resources"
}

data "azurerm_api_management_api" "example" {
  name                = "search-api"
  api_management_name = data.azurerm_api_management.example.name
  resource_group_name = data.azurerm_api_management.example.resource_group_name
  revision            = "2"
}

data "azurerm_api_management_gateway" "example" {
  name              = "example-gateway"
  api_management_id = data.azurerm_api_management.example.id
}

resource "azurerm_api_management_gateway_api" "example" {
  gateway_id = data.azurerm_api_management_gateway.example.id
  api_id     = data.azurerm_api_management_api.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `gateway_id` - (Required) The Identifier for the API Management Gateway. Changing this forces a new API Management Gateway API to be created.

* `api_id` - (Required) The Identifier of the API Management API within the API Management Service. Changing this forces a new API Management Gateway API to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the API Management Gateway API.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Management Gateway API.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management Gateway API.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Management Gateway API.

## Import

API Management Gateway APIs can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_gateway_api.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.ApiManagement/service/service1/gateways/gateway1/apis/api1
```
