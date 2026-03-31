---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_relay_namespace_authorization_rule"
description: |-
  Gets information about an existing Azure Relay Namespace Authorization Rule.
---

# Data Source: azurerm_relay_namespace_authorization_rule

Use this data source to access information about an existing Azure Relay Namespace Authorization Rule.

## Example Usage

```hcl
data "azurerm_relay_namespace_authorization_rule" "example" {
  name                = "example"
  resource_group_name = "example-rg"
  namespace_name      = "example-relay"
}

output "id" {
  value = data.azurerm_relay_namespace_authorization_rule.example.id
}

output "primary_key" {
  value = data.azurerm_relay_namespace_authorization_rule.example.primary_key
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Azure Relay Namespace Authorization Rule.

* `namespace_name` - (Required) Name of the Azure Relay Namespace for which this Azure Relay Namespace Authorization Rule exists.

* `resource_group_name` - (Required) The name of the Resource Group where the Azure Relay Namespace Authorization Rule exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Azure Relay Namespace Authorization Rule.

* `listen` - Whether listen access is granted to this Azure Relay Namespace Authorization Rule.

* `manage` - Whether manage access is granted to this Azure Relay Namespace Authorization Rule.

* `primary_key` - The Primary Key for the Azure Relay Namespace Authorization Rule.

* `primary_connection_string` - The Primary Connection String for the Azure Relay Namespace Authorization Rule.

* `secondary_key` - The Secondary Key for the Azure Relay Namespace Authorization Rule.

* `secondary_connection_string` - The Secondary Connection String for the Azure Relay Namespace Authorization Rule.

* `send` - Whether send access is granted to this Azure Relay Namespace Authorization Rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Azure Relay Namespace Authorization Rule.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.Relay` - 2021-11-01
