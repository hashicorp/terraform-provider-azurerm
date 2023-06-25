---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_api_management_gateway_host_name_configuration"
description: |-
  Gets information about an existing API Management Gateway Host Configuration.
---

# Data Source: azurerm_api_management_gateway_host_name_configuration

Use this data source to access information about an existing API Management Gateway Host Configuration.

## Example Usage

```hcl
data "azurerm_api_management" "example" {
  name                = "example-apim"
  resource_group_name = "example-resources"
}

data "azurerm_api_management_gateway" "example" {
  name              = "example-gateway"
  api_management_id = data.azurerm_api_management.main.id
}

data "azurerm_api_management_gateway_host_name_configuration" "example" {
  name              = "example-host-configuration"
  api_management_id = data.azurerm_api_management.example.id
  gateway_name      = data.azurerm_api_management_gateway.example.name
}

output "host_name" {
  value = data.azurerm_api_management_gateway_host_name_configuration.example.host_name
}
```

## Arguments Reference

The following arguments are supported:

* `api_management_id` - (Required) The ID of the API Management Service.

* `gateway_name` - (Required) The name of the API Management Gateway.
* 
* `name` - (Required) The name of the API Management Gateway Host Name Configuration.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `certificate_id` - The ID of the certificate used for TLS connection establishment.

* `host_name` - The host name used for the API Management Gateway Host Name Configuration.

* `http2_enabled` - Whether HTTP/2.0 is supported.

* `id` - The ID of the API Management Gateway Host Configuration.

* `request_client_certificate_enabled` - Whether the API Management Gateway requests a client certificate.

* `tls10_enabled` - Whether TLS 1.0 is supported.

* `tls11_enabled` - Whether TLS 1.1 is supported.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the API Management Gateway Host Configuration.
