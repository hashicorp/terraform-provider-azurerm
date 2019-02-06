---
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
data "azurerm_servicebus_namespace" "test" {
  name                = "examplenamespace"
  resource_group_name = "example-resources"
}

output "location" {
  value = "${data.azurerm_servicebus_namespace.test.location}"
}
```

## Argument Reference

* `name` - (Required) Specifies the name of the ServiceBus Namespace.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the ServiceBus Namespace exists.

## Attributes Reference

* `location` - The location of the resource group.
