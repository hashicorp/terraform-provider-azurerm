---
subcategory: "Container Apps"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_container_app_environment_maintenance_configuration"
description: |-
    Lists Container App Environment Maintenance Configuration resources.
---

# List resource: azurerm_container_app_environment_maintenance_configuration

Lists Container App Environment Maintenance Configuration resources.

## Example Usage

### List Maintenance Configurations in a Container App Environment

```hcl
list "azurerm_container_app_environment_maintenance_configuration" "example" {
  provider = azurerm
  config {
    container_app_environment_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.App/managedEnvironments/myEnvironment"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `container_app_environment_id` - (Required) The ID of the Container App Environment to query.
