---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_gateway"
description: |-
  Gets information about an existing API Management Gateway.
---

# Data Source: azurerm_api_management_gateway

Use this data source to access information about an existing API Management Gateway.

## Example Usage

```hcl
data "azurerm_api_management" "example" {
  name                = "example-apim"
  resource_group_name = "example-rg"
}

data "azurerm_api_management_gateway" "example" {
  name              = "example-api-gateway"
  api_management_id = data.azurerm_api_management.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - The name of the API Management Gateway.

* `api_management_id` - The ID of the API Management Service in which the Gateway exists.

## Attributes Reference

* `id` - The ID of the API Management Gateway.

* `location` - A `location` block as documented below.

* `description` - The description of the API Management Gateway.

---

A `location` block exports the following:

* `name` - A canonical name for the geographic or physical location.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the API Management Gateway.
