---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_eventhub_namespace"
sidebar_current: "docs-azurerm-resource-messaging-eventhub-namespace-x"
description: |-
  Manages an EventHub Namespace.
---

# azurerm_eventhub_namespace

Manages an EventHub Namespace.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "resourceGroup1"
  location = "West US"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acceptanceTestEventHubNamespace"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
  capacity            = 2

  tags = {
    environment = "Production"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the EventHub Namespace resource. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the namespace. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `sku` - (Required) Defines which tier to use. Valid options are `Basic` and `Standard`.

* `capacity` - (Optional) Specifies the Capacity / Throughput Units for a `Standard` SKU namespace. Valid values range from 1 - 20.

* `auto_inflate_enabled` - (Optional) Is Auto Inflate enabled for the EventHub Namespace?

* `maximum_throughput_units` - (Optional) Specifies the maximum number of throughput units when Auto Inflate is Enabled. Valid values range from 1 - 20.

* `kafka_enabled` - (Optional) Is Kafka enabled for the EventHub Namespace? Defaults to `false`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The EventHub Namespace ID.

The following attributes are exported only if there is an authorization rule named
`RootManageSharedAccessKey` which is created automatically by Azure.

* `default_primary_connection_string` - The primary connection string for the authorization
    rule `RootManageSharedAccessKey`.

* `default_secondary_connection_string` - The secondary connection string for the
    authorization rule `RootManageSharedAccessKey`.

* `default_primary_key` - The primary access key for the authorization rule `RootManageSharedAccessKey`.

* `default_secondary_key` - The secondary access key for the authorization rule `RootManageSharedAccessKey`.

## Import

EventHub Namespaces can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_eventhub_namespace.namespace1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.EventHub/namespaces/namespace1
```
