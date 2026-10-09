---
subcategory: "Container Apps"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_container_app_connection"
description: |-
  Lists Service Connectors for a Container App.
---

# List resource: azurerm_container_app_connection

Lists Service Connectors for a Container App.

## Example Usage

### List all connections for a specific Container App

```hcl
list "azurerm_container_app_connection" "example" {
  provider = azurerm
  config {
    container_app_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.App/containerApps/containerApp1"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `container_app_id` - (Required) The ID of the Container App whose connections should be listed.

~> **Note:** Authentication secrets are not returned by the API. When using `include_resource = true`, any required secrets must be supplied separately before applying generated configuration.
