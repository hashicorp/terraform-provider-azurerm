---
subcategory: "Sentinel"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_sentinel_data_connector_office_advanced_threat_protection"
description: |-
  Manages a Office ATP Data Connector.
---

# azurerm_sentinel_data_connector_office_advanced_threat_protection

Manages a Office ATP Data Connector.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "West Europe"
}

resource "azurerm_log_analytics_workspace" "example" {
  name                = "example-workspace"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "PerGB2018"
}

resource "azurerm_log_analytics_solution" "example" {
  solution_name         = "SecurityInsights"
  location              = azurerm_resource_group.example.location
  resource_group_name   = azurerm_resource_group.example.name
  workspace_resource_id = azurerm_log_analytics_workspace.example.id
  workspace_name        = azurerm_log_analytics_workspace.example.name

  plan {
    publisher = "Microsoft"
    product   = "OMSGallery/SecurityInsights"
  }
}

resource "azurerm_sentinel_data_connector_office_advanced_threat_protection" "example" {
  name                       = "example"
  log_analytics_workspace_id = "TODO"
}

resource "azurerm_sentinel_data_connector_office_advanced_threat_protection" "example" {
  name                       = "example"
  log_analytics_workspace_id = azurerm_log_analytics_solution.example.workspace_resource_id
  depends_on                 = [azurerm_log_analytics_solution.example]
}
```

## Arguments Reference

The following arguments are supported:

- `name` - (Required) The name which should be used for this Office ATP Data Connector. Changing this forces a new Office ATP Data Connector to be created.

- `log_analytics_workspace_id` - (Required) The ID of the Log Analytics Workspace that this Office ATP Data Connector resides in. Changing this forces a new Office ATP Data Connector to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

- `id` - The ID of the Office ATP Data Connector.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

- `create` - (Defaults to 30 minutes) Used when creating the Office ATP Data Connector.
- `read` - (Defaults to 5 minutes) Used when retrieving the Office ATP Data Connector.
- `update` - (Defaults to 30 minutes) Used when updating the Office ATP Data Connector.
- `delete` - (Defaults to 30 minutes) Used when deleting the Office ATP Data Connector.

## Import

Office ATP Data Connectors can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_sentinel_data_connector_office_advanced_threat_protection.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.OperationalInsights/workspaces/workspace1/providers/Microsoft.SecurityInsights/dataConnectors/dc1
```
