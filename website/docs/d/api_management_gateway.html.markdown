---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_api_management_gateway"
description: |-
  Gets information about an existing API Management Gateway.
---

# Data Source: azurerm_api_management_gateway

Use this data source to access information about an existing API Management Gateway.

## Example Usage

```hcl
data "azurerm_api_management_gateway" "example" {
  resource_group_name = "existing"
  gateway_id = "my-gateway"
  api_management_name = "existing"
}

output "id" {
  value = data.azurerm_api_management_gateway.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `api_management_name` - The Name of the API Management Service in which this Gateway exists.

* `gateway_id` - The Identifier for the API Management Gateway.

* `resource_group_name` - (Required) The name of the Resource Group where the API Management Gateway exists. Changing this forces a new API Management Gateway to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `description` - The Description of the API Management Gateway.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the API Management Gateway.