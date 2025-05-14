---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_servicebus_namespace_disaster_recovery_config"
description: |-
  Gets information about an existing Service Bus Disaster Recovery Config.
---

# Data Source: azurerm_servicebus_namespace_disaster_recovery_config

Use this data source to access information about an existing Service Bus Disaster Recovery Config.

## Example Usage

```hcl
data "azurerm_servicebus_namespace_disaster_recovery_config" "example" {
  name         = "existing"
  namespace_id = "example-namespace-id"
}

output "id" {
  value = data.azurerm_servicebus_namespace_disaster_recovery_config.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Service Bus Disaster Recovery Config.

* `namespace_id` - (Required) The ID of the Service Bus Namespace.

---

* `alias_authorization_rule_id` - (Optional) The Shared access policies used to access the connection string for the alias.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Service Bus Disaster Recovery Config.

* `default_primary_key` - The primary access key for the authorization rule `RootManageSharedAccessKey`.

* `default_secondary_key` - The secondary access key for the authorization rule `RootManageSharedAccessKey`.

* `partner_namespace_id` - The ID of the Service Bus Namespace to replicate to.

* `primary_connection_string_alias` - The alias Primary Connection String for the ServiceBus Namespace.

* `secondary_connection_string_alias` - The alias Secondary Connection String for the ServiceBus Namespace

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Service Bus Disaster Recovery Config.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.ServiceBus`: 2021-06-01-preview
