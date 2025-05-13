---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_servicebus_namespace_disaster_recovery_config"
description: |-
  Manages a Disaster Recovery Config for a Service Bus Namespace.
---

# azurerm_servicebus_namespace_disaster_recovery_config

Manages a Disaster Recovery Config for a Service Bus Namespace.

~> **Note:** Disaster Recovery Config is a Premium SKU only capability.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "servicebus-replication"
  location = "West Europe"
}

resource "azurerm_servicebus_namespace" "primary" {
  name                = "servicebus-primary"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "Premium"
  capacity            = "1"
}

resource "azurerm_servicebus_namespace" "secondary" {
  name                = "servicebus-secondary"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "Premium"
  capacity            = "1"
}

resource "azurerm_servicebus_namespace_authorization_rule" "example" {
  name         = "examplerule"
  namespace_id = azurerm_servicebus_namespace.example.id

  listen = true
  send   = true
  manage = false
}

resource "azurerm_servicebus_namespace_disaster_recovery_config" "example" {
  name                        = "servicebus-alias-name"
  primary_namespace_id        = azurerm_servicebus_namespace.primary.id
  partner_namespace_id        = azurerm_servicebus_namespace.secondary.id
  alias_authorization_rule_id = azurerm_servicebus_namespace_authorization_rule.example.id
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Disaster Recovery Config. This is the alias DNS name that will be created. Changing this forces a new resource to be created.

* `primary_namespace_id` - (Required) The ID of the primary Service Bus Namespace to replicate. Changing this forces a new resource to be created.

* `partner_namespace_id` - (Required) The ID of the Service Bus Namespace to replicate to.

* `alias_authorization_rule_id` - (Optional) The Shared access policies used to access the connection string for the alias.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The Service Bus Namespace Disaster Recovery Config ID.

* `primary_connection_string_alias` - The alias Primary Connection String for the ServiceBus Namespace.

* `secondary_connection_string_alias` - The alias Secondary Connection String for the ServiceBus Namespace

* `default_primary_key` - The primary access key for the authorization rule `RootManageSharedAccessKey`.

* `default_secondary_key` - The secondary access key for the authorization rule `RootManageSharedAccessKey`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Service Bus Namespace Disaster Recovery Config.
* `read` - (Defaults to 5 minutes) Used when retrieving the Service Bus Namespace Disaster Recovery Config.
* `update` - (Defaults to 30 minutes) Used when updating the Service Bus Namespace Disaster Recovery Config.
* `delete` - (Defaults to 30 minutes) Used when deleting the Service Bus Namespace Disaster Recovery Config.

## Import

Service Bus DR configs can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_servicebus_namespace_disaster_recovery_config.config1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ServiceBus/namespaces/namespace1/disasterRecoveryConfigs/config1
```
