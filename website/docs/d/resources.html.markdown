---
subcategory: "Base"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_resources"
sidebar_current: "docs-azurerm-datasource-resources"
description: |-
  Gets information about an existing Resources.
---

# Data Source: azurerm_resources

Use this data source to access information about existing resources.

## Example Usage

```hcl
# Get Resources from a Resource Group
data "azurerm_resources" "example" {
  resource_group_name = "example-resources"
}

# Get Resources with specific Tags
data "azurerm_resources" "example" {
  resource_group_name = "example-resources"

  required_tags {
    environment = "production"
    role        = "webserver"
  }
}
```

## Argument Reference

~> **NOTE:** At least one of `name`, `resource_group_name` or `type` must be specified.

* `name` - (Optional) The name of the Resource.

* `resource_group_name` - (Optional) The name of the Resource group where the Resources are located.

* `type` - (Optional) The Resource Type of the Resources you want to list (e.g. `Microsoft.Network/virtualNetworks`). A full list of available Resource Types can be found [here](https://docs.microsoft.com/en-us/azure/azure-resource-manager/azure-services-resource-providers).

* `required_tags` - (Optional) A mapping of tags which the resource has to have in order to be included in the result.

## Attributes Reference

* `resources` - One or more `resource` blocks as defined below.

---

The `resource` block contains:

* `name` - The name of this resource.

* `id` - The Resource ID of this resource.

* `type` - The type of this resoource.

* `location` - The location of this resource.

* `tags` - The type of resource that this is, such as `Microsoft.Network/virtualNetworks`.
