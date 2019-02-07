---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_availability_set"
sidebar_current: "docs-azurerm-datasource-availability_set"
description: |-
  Gets information about an existing Availability Set.
---

# Data Source: azurerm_availability_set

Use this data source to access information about an existing Availability Set.

## Example Usage

```hcl
data "azurerm_availability_set" "test" {
  name                = "tf-appsecuritygroup"
  resource_group_name = "my-resource-group"
}

output "availability_set_id" {
  value = "${data.azurerm_availability_set.test.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - The name of the Availability Set.

* `resource_group_name` - The name of the resource group in which the Availability Set exists.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Availability Set.

* `location` - The supported Azure location where the Availability Set exists.

* `tags` - A mapping of tags assigned to the resource.
