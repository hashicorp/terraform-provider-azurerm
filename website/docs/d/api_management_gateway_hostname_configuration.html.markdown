---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_api_management_gateway_hostname_configuration"
description: |-
  Gets information about an existing API Management Gateway Host Configuration.
---

# Data Source: azurerm_api_management_gateway_hostname_configuration

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

data "azurerm_api_management_gateway_hostname_configuration" "example" {
  name = "example-host-configuration"
  api_management_gateway_id = data.azurerm_api_management_gateway.example.id
}

output "id" {
  value = data.azurerm_api_management_gateway_hostname_configuration.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `api_management_gateway_id` - (Required) The ID of the Gateway.

* `name` - (Required) The name of this API Management Gateway Host Configuration.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the API Management Gateway Host Configuration.

* `certificate_id` - The ID of the certificate.

* `hostname` - The hostname.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the API Management Gateway Host Configuration.
