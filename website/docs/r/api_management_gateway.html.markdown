---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_gateway"
description: |-
  Manages an API Management Gateway.
---

# azurerm_api_management_gateway

Manages an API Management Gateway.

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
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Consumption_0"
}

resource "azurerm_api_management_gateway" "example" {
  name              = "example-gateway"
  api_management_id = azurerm_api_management.example.id
  description       = "Example API Management gateway"

  location_data {
    name     = "example name"
    city     = "example city"
    district = "example district"
    region   = "example region"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for the API Management Gateway. Changing this forces a new API Management Gateway to be created.

* `api_management_name` - (Required) The name of the API Management Service in which the gateway will be created. Changing this forces a new API Management Gateway resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which the API Management Gateway exists.

* `location_data` - (Required) A `location_data` block as documented below.

* `description` - (Optional) The description of the API Management Gateway.

---

A `location_data` block supports the following:

* `name` - (Required) A canonical name for the geographic or physical location.

* `city` - (Optional) The city or locality where the resource is located.

* `district` - (Optional) The district, state, or province where the resource is located.

* `region` - (Optional) The country or region where the resource is located.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the API Management Gateway.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Management Gateway.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management Gateway.
* `update` - (Defaults to 30 minutes) Used when updating the API Management Gateway.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Management Gateway.

## Import

API Management Gateways can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_gateway.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ApiManagement/service/service1/gateways/gateway1
```
