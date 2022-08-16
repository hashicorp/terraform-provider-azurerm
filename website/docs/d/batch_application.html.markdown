---
subcategory: "Batch"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_batch_application"
description: |-
  Get information about an existing Batch Application instance
---

# Data Source: azurerm_batch_application

Use this data source to access information about an existing Batch Application instance.

## Example Usage

```hcl
data "azurerm_batch_application" "example" {
  name                = "testapplication"
  resource_group_name = "test"
  account_name        = "testbatchaccount"
}

output "batch_application_id" {
  value = data.azurerm_batch_application.example.id
}
```

## Argument Reference

* `name` - The name of the Application.

* `resource_group_name` - The name of the Resource Group where this Batch account exists.

* `account_name` - The name of the Batch account.

## Attributes Reference

The following attributes are exported:

* `id` - The Batch application ID.

* `name` - The Batch application name.

* `allow_updates` - May packages within the application be overwritten using the same version string.

* `default_version` - The package to use if a client requests the application but does not specify a version.

* `display_name` - The display name for the application.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Batch Application.
