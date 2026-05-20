---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_relay_hybrid_connection"
description: |-
  Gets information about an existing Azure Relay Hybrid Connection.
---

# Data Source: azurerm_relay_hybrid_connection

Use this data source to access information about an existing Azure Relay Hybrid Connection.

## Example Usage
z
```hcl
data "azurerm_relay_hybrid_connection" "example" {
  name                 = "example"
  resource_group_name  = "example-rg"
  relay_namespace_name = "example-relay"
}

output "id" {
  value = data.azurerm_relay_hybrid_connection.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Azure Relay Hybrid Connection.

* `relay_namespace_name` - (Required) The name of the Azure Relay in which the Azure Relay Hybrid Connection exists.

* `resource_group_name` - (Required) The name of the Resource Group where the Azure Relay Hybrid Connection exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Azure Relay Hybrid Connection.

* `requires_client_authorization` - Specifies if client authorization is needed for this Azure Relay Hybrid Connection.

* `user_metadata` - The usermetadata is a placeholder to store user-defined string data for the Azure Relay Hybrid Connection endpoint.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Azure Relay Hybrid Connection.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.Relay` - 2021-11-01
