---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_resource"
sidebar_current: "docs-azurerm-datasource-resource"
description: |-
  Gets information about an existing Resource.
---

# Data Source: azurerm_resource

Use this data source to access information about an existing Resource.

## Example Usage

```hcl
// Get Resource by Resource ID
data "azurerm_resource" "test" {
  resource_id = "/subscriptions/00000000-0000-0000-000000000000/resourceGroups/myResourceGroup/providers/Microsoft.Storage/storageAccounts/myStorageAcc"
}

// Get Resources from a Resource Group
data "azurerm_resource" "test" {
  resource_group_name = "myResourceGroup"
}

// Get Resources with specific Tags
data "azurerm_resource" "test" {
  required_tags {
    environment = "production"
    role        = "webserver"
  }
}
```

## Argument Reference

* `resource_id` - (Optional) The fully qualified ID of the resource. `resource_id` can't be used with the other Arguments.

* `name` - (Optional) The name of the Resource.

* `resource_group_name` - (Optional) The name of the Resource group where the Resources are located.

* `type` - (Optional) The Resource Type of the Resources you want to list (e.g. `Microsoft.Network/virtualNetworks`). A full list of available Resource Types can be found [here](https://docs.microsoft.com/en-us/azure/azure-resource-manager/azure-services-resource-providers).

* `required_tags` - (Optional) A mapping of tags which the Resource has to have in order to be included in the result.

## Attributes Reference

* `resources` - One or more `resource` blocks as defined below.

---

The `resource` block contains:

* `name` - The location of the resource group.

* `id` - A mapping of tags assigned to the resource group.

* `type` - A mapping of tags assigned to the resource group.

* `location` - A mapping of tags assigned to the resource group.

* `tags` - A mapping of tags assigned to the resource group.
