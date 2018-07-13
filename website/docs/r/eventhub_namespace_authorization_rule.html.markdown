---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_eventhub_namespace_authorization_rule"
sidebar_current: "docs-azurerm-resource-eventhub-namespace-authorization-rule"
description: |-
  Manages a Event Hub Namespace authorization Rule within an Event Hub.
---

# azurerm_eventhub_namespace_authorization_rule

Manages a Event Hub Namespace authorization Rule within an Event Hub.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "resourceGroup1"
  location = "West US"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acceptanceTestEventHubNamespace"
  location            = "West US"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Basic"
  capacity            = 2

  tags {
    environment = "Production"
  }
}

resource "azurerm_eventhub_namespace_authorization_rule" "test" {
  name                = "navi"
  namespace_name      = "${azurerm_eventhub_namespace.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  
  listen              = true
  send                = false
  manage              = false
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the EventHub Namespace Authorization Rule resource. Changing this forces a new resource to be created.

* `namespace_name` - (Required) Specifies the name of the EventHub Namespace. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the EventHub Namespace exists. Changing this forces a new resource to be created.

~> **NOTE** At least one of the 3 permissions below needs to be set.

* `listen` - (Optional) Grants listen access to this this Authorization Rule. Defaults to `false`.

* `send` - (Optional) Grants send access to this this Authorization Rule. Defaults to `false`.

* `manage` - (Optional) Grants manage access to this this Authorization Rule. When this property is `true` - both `listen` and `send` must be too. Defaults to `false`.

## Attributes Reference

The following attributes are exported:

* `id` - The EventHub ID.

* `primary_key` - The Primary Key for the Event Hubs authorization Rule.

* `primary_connection_string` - The Primary Connection String for the Event Hubs authorization Rule.

* `secondary_key` - The Secondary Key for the Event Hubs authorization Rule.

* `secondary_connection_string` - The Secondary Connection String for the Event Hubs authorization Rule.

## Import

EventHubs can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_eventhub_namespace_authorization_rule.rule1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.EventHub/namespaces/namespace1/authorizationRules/rule1
```
