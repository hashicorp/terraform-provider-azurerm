---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_gateway_certificate_authority"
description: |-
  Manages an API Management Gateway Certificate Authority.
---

# azurerm_api_management_gateway_certificate_authority

Manages an API Management Gateway Certificate Authority.

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

resource "azurerm_api_management_certificate" "example" {
  name                = "example-cert"
  api_management_name = azurerm_api_management.example.name
  resource_group_name = azurerm_resource_group.example.name
  data                = filebase64("example.pfx")
}

resource "azurerm_api_management_gateway_certificate_authority" "example" {
  api_management_id = azurerm_api_management.example.id
  certificate_name  = azurerm_api_management_certificate.example.name
  gateway_name      = azurerm_api_management_gateway.example.name
  is_trusted        = true
}
```

## Arguments Reference

The following arguments are supported:

* `api_management_id` - (Required) The ID of the API Management Service. Changing this forces a new resource to be created.

* `certificate_name` - (Required) The name of the API Management Certificate. Changing this forces a new resource to be created.

* `gateway_name` - (Required) The name of the API Management Gateway. Changing this forces a new resource to be created.

* `is_trusted` - (Optional) Whether the API Management Gateway Certificate Authority is trusted.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the API Management Gateway Certificate Authority.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Management Gateway Certificate Authority.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management Gateway Certificate Authority.
* `update` - (Defaults to 30 minutes) Used when updating the API Management Gateway Certificate Authority.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Management Gateway Certificate Authority.

## Import

API Management Gateway Certificate Authority can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_gateway_certificate_authority.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ApiManagement/service/service1/gateways/gateway1/certificateAuthorities/cert1
```
