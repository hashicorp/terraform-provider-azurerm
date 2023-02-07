---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_gateway_host_name_configuration"
description: |-
  Manages an API Management Gateway Host Name Configuration.
---

# azurerm_api_management_gateway_host_name_configuration

Manages an API Management Gateway Host Name Configuration.

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

resource "azurerm_api_management_gateway_host_name_configuration" "example" {
  name              = "example-host-name-configuration"
  api_management_id = azurerm_api_management.example.id
  gateway_name      = azurerm_api_management_gateway.example.name

  certificate_id                     = azurerm_api_management_certificate.example.id
  host_name                          = "example-host-name"
  request_client_certificate_enabled = true
  http2_enabled                      = true
  tls10_enabled                      = true
  tls11_enabled                      = false
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the API Management Gateway Host Name Configuration. Changing this forces a new resource to be created.

* `api_management_id` - (Required) The ID of the API Management Service. Changing this forces a new resource to be created.

* `gateway_name` - (Required) The name of the API Management Gateway. Changing this forces a new resource to be created.

* `certificate_id` - (Required) The certificate ID to be used for TLS connection establishment.

* `host_name` - (Required) The host name to use for the API Management Gateway Host Name Configuration.

* `request_client_certificate_enabled` - (Optional) Whether the API Management Gateway requests a client certificate.

* `http2_enabled` - (Optional) Whether HTTP/2.0 is supported. Defaults to `true`.

* `tls10_enabled` - (Optional) Whether TLS 1.0 is supported.

* `tls11_enabled` - (Optional) Whether TLS 1.1 is supported.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the API Management Gateway Host Name Configuration.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Management Gateway Host Name Configuration.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management Gateway Host Name Configuration.
* `update` - (Defaults to 30 minutes) Used when updating the API Management Gateway Host Name Configuration.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Management Gateway Host Name Configuration.

## Import

API Management Gateway Host Name Configuration can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_gateway_host_name_configuration.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ApiManagement/service/service1/gateways/gateway1/hostnameConfigurations/hc1
```
