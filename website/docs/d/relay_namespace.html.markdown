---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_relay_namespace"
description: |-
  Gets information about an existing Relay Namespace.
---

# Data Source: azurerm_relay_namespace

Use this data source to access information about an existing Relay Namespace.

## Example Usage

```hcl
data "azurerm_relay_namespace" "example" {
  name                = "existing"
  resource_group_name = "existing"
}

output "id" {
  value = data.azurerm_relay_namespace.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Relay Namespace.

* `resource_group_name` - (Required) The name of the Resource Group where the Relay Namespace exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Relay Namespace.

* `location` - The Azure Region where the Relay Namespace exists.

* `metric_id` - The Identifier for Azure Insights metrics.

-> **Note:** The `primary_connection_string`, `primary_key`, `secondary_connection_string` and `secondary_key` attributes are only populated when the `RootManageSharedAccessKey` authorization rule exists and the credentials used by Terraform are allowed to list its keys.

* `primary_connection_string` - The primary connection string for the authorization rule `RootManageSharedAccessKey`.

* `primary_key` - The primary access key for the authorization rule `RootManageSharedAccessKey`.

* `secondary_connection_string` - The secondary connection string for the authorization rule `RootManageSharedAccessKey`.

* `secondary_key` - The secondary access key for the authorization rule `RootManageSharedAccessKey`.

* `sku_name` - The name of the SKU used by this Relay Namespace.

* `tags` - A mapping of tags assigned to the Relay Namespace.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Relay Namespace.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.Relay` - 2021-11-01
