---
subcategory: "Connections"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_managed_api"
description: |-
    Gets information about an existing Managed API.
---

# Data Source: azurerm_managed_api

Uses this data source to access information about an existing Managed API.

## Example Usage

```hcl
data "azurerm_managed_api" "example" {
  name     = "servicebus"
  location = "West Europe"
}

output "id" {
  value = data.azurerm_managed_api.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Managed API.

* `location` - (Required) The Azure location for this Managed API.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Managed API.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Managed API.
