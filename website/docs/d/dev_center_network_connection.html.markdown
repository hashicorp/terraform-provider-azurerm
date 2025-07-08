---
subcategory: "Dev Center"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_dev_center_network_connection"
description: |-
  Gets information about an existing Dev Center Network Connection.
---

# Data Source: azurerm_dev_center_network_connection

Use this data source to access information about an existing Dev Center Network Connection.

## Example Usage

```hcl
data "azurerm_dev_center_network_connection" "example" {
  name                = "example"
  resource_group_name = "example-resource-group"
}

output "id" {
  value = data.azurerm_dev_center_network_connection.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Dev Center Network Connection.

* `resource_group_name` - (Required) The name of the Resource Group where the Dev Center Network Connection exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Dev Center Network Connection.

* `domain_join_type` - The Azure Active Directory Join type.

* `domain_name` - The name of the Azure Active Directory domain.

* `domain_username` - The username of the Azure Active Directory account (user or service account) that has permissions to create computer objects in Active Directory.

* `location` - The Azure Region where the Dev Center Network Connection exists.

* `organization_unit` - The Azure Active Directory domain Organization Unit (OU).

* `subnet_id` - The ID of the Subnet that is used to attach Virtual Machines.

* `tags` - A mapping of tags assigned to the Dev Center Network Connection.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Dev Center Network Connection.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.DevCenter`: 2025-02-01
