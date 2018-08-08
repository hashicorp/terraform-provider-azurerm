---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_resource_group"
sidebar_current: "docs-azurerm-datasource-resource-group"
description: |-
  Get information about the specified resource group.
---

# Data Source: azurerm_resource_group

Use this data source to access the properties of an Azure resource group.

## Example Usage

```hcl
data "azurerm_resource_group" "example" {
  name = "example-resources"
}

output "location" {
  value = "${data.azurerm_resource_group.example.location}"
}
```

## Argument Reference

* `name` - (Required) Specifies the name of the resource group.

~> **NOTE:** If the specified location doesn't match the actual resource group location, an error message with the actual location value will be shown.

## Attributes Reference

* `location` - The location of the resource group.
* `tags` - A mapping of tags assigned to the resource group.
