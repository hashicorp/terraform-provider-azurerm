---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_gateway_hostname_configuration"
description: |-
  Manages an API Management Gateway Host Configuration.
---

# azurerm_api_management_gateway_hostname_configuration

Manages an API Management Gateway Host Configuration.

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
}

resource "azurerm_api_management_certificate" "example" {
  name                = "example-cert"
  api_management_name = azurerm_api_management.example.name
  resource_group_name = azurerm_resource_group.example.name
  data                = filebase64("data/cert.pfx")
}

resource "azurerm_api_management_gateway_hostname_configuration" "example" {
  name                      = "example-hostname-configuration"
  api_management_gateway_id = azurerm_api_management_gateway.example.id
  hostname                  = "app.example.com"
  certificate_id            = azurerm_api_management_certificate.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `api_management_gateway_id` - (Required) The ID of the Gateway. Changing this forces a new API Management Gateway Host Configuration to be created.

* `certificate_id` - (Required) The ID of the certificate.

* `hostname` - (Required) The hostname.

* `name` - (Required) The name which should be used for this API Management Gateway Host Configuration. Changing this forces a new API Management Gateway Host Configuration to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the API Management Gateway Host Configuration.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Management Gateway Host Configuration.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management Gateway Host Configuration.
* `update` - (Defaults to 30 minutes) Used when updating the API Management Gateway Host Configuration.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Management Gateway Host Configuration.

## Import

API Management Gateway Host Configurations can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_gateway_hostname_configuration.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ApiManagement/service/service1/gateways/gateway1/hostConfigurations/hostConfiguration1
```
