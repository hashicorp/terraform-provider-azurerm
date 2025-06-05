---
subcategory: "Managed Applications"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_managed_application_definition"
description: |-
  Gets information about an existing Managed Application Definition
---

# Data Source: azurerm_managed_application_definition

Uses this data source to access information about an existing Managed Application Definition.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

data "azurerm_managed_application_definition" "example" {
  name                = "examplemanagedappdef"
  resource_group_name = "exampleresources"
}

output "id" {
  value = data.azurerm_managed_application_definition.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Managed Application Definition.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where this Managed Application Definition exists.

* `location` - The Azure location where the resource exists.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Managed Application Definition.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Managed Application Definition.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.Solutions`: 2021-07-01
