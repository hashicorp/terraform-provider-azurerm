---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_servicebus_namespace"
sidebar_current: "docs-azurerm-datasource-servicebus-namespace"
description: |-
  Gets information about an existing ServiceBus Namespace.
---

# Data Source: azurerm_servicebus_namespace

Use this data source to access information about an existing ServiceBus Namespace.

## Example Usage

```hcl
data "azurerm_servicebus_namespace" "example" {
  name                = "examplenamespace"
  resource_group_name = "example-resources"
}

output "location" {
  value = "${data.azurerm_servicebus_namespace.example.location}"
}
```

## Argument Reference

* `name` - (Required) Specifies the name of the ServiceBus Namespace.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the ServiceBus Namespace exists.

## Attributes Reference

* `location` - The location of the Resource Group in which the ServiceBus Namespace exists.

* `sku` - The Tier used for the ServiceBus Namespace.

* `capacity` - The capacity of the ServiceBus Namespace.

* `zone_redundant` - Whether or not this ServiceBus Namespace is zone redundant.

* `tags` - A mapping of tags assigned to the resource.

The following attributes are exported only if there is an authorization rule named
`RootManageSharedAccessKey` which is created automatically by Azure.

* `default_primary_connection_string` - The primary connection string for the authorization
    rule `RootManageSharedAccessKey`.

* `default_secondary_connection_string` - The secondary connection string for the
    authorization rule `RootManageSharedAccessKey`.

* `default_primary_key` - The primary access key for the authorization rule `RootManageSharedAccessKey`.

* `default_secondary_key` - The secondary access key for the authorization rule `RootManageSharedAccessKey`.
