---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_eventhub_namespace_disaster_recovery_config"
description: |-
  Manages an Disaster Recovery Config for an Event Hub Namespace.
---

# azurerm_eventhub_namespace_disaster_recovery_config

Manages an Disaster Recovery Config for an Event Hub Namespace.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "eventhub-replication"
  location = "West Europe"
}

resource "azurerm_eventhub_namespace" "primary" {
  name                = "eventhub-primary"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "Standard"
}

resource "azurerm_eventhub_namespace" "secondary" {
  name                = "eventhub-secondary"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "Standard"
}

resource "azurerm_eventhub_namespace_disaster_recovery_config" "example" {
  name                 = "replicate-eventhub"
  resource_group_name  = azurerm_resource_group.example.name
  namespace_name       = azurerm_eventhub_namespace.primary.name
  partner_namespace_id = azurerm_eventhub_namespace.secondary.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Disaster Recovery Config. Changing this forces a new resource to be created.

* `namespace_name` - (Required) Specifies the name of the primary EventHub Namespace to replicate. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the Disaster Recovery Config exists. Changing this forces a new resource to be created.

* `partner_namespace_id` - (Required) The ID of the EventHub Namespace to replicate to.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The EventHub Namespace Disaster Recovery Config ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the EventHub Namespace Disaster Recovery Config.
* `read` - (Defaults to 5 minutes) Used when retrieving the EventHub Namespace Disaster Recovery Config.
* `update` - (Defaults to 30 minutes) Used when updating the EventHub Namespace Disaster Recovery Config.
* `delete` - (Defaults to 30 minutes) Used when deleting the EventHub Namespace Disaster Recovery Config.

## Import

EventHubs can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_eventhub_namespace_disaster_recovery_config.config1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.EventHub/namespaces/namespace1/disasterRecoveryConfigs/config1
```
