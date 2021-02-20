---
subcategory: "Sentinel"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_sentinel_data_connector_microsoft_cloud_app_security"
description: |-
  Manages a Microsoft Cloud App Security Data Connector.
---

# azurerm_sentinel_data_connector_microsoft_cloud_app_security

Manages a Microsoft Cloud App Security Data Connector.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "west europe"
}

resource "azurerm_log_analytics_workspace" "example" {
  name                = "example-workspace"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "PerGB2018"
}

resource "azurerm_sentinel_data_connector_microsoft_cloud_app_security" "example" {
  name                       = "example"
  log_analytics_workspace_id = azurerm_log_analytics_workspace.example.id
}
```

## Arguments Reference

The following arguments are supported:

- `log_analytics_workspace_id` - (Required) The ID of the Log Analytics Workspace that this Microsoft Cloud App Security Data Connector resides in. Changing this forces a new Microsoft Cloud App Security Data Connector to be created.

* `name` - (Required) The name which should be used for this Microsoft Cloud App Security Data Connector. Changing this forces a new Microsoft Cloud App Security Data Connector to be created.

---

- `alerts_enabled` - (Optional) Should the alerts be enabled? Defaults to `true`.

- `discovery_logs_enabled` - (Optional) Should the TODO be enabled? Defaults to `true`.

-> **NOTE**: At least one of `alerts_enabled` and `discovery_logs_enabled` has to be specified.

- `tenant_id` - (Optional) The ID of the tenant that this Microsoft Cloud App Security Data Connector connects to. Changing this forces a new Microsoft Cloud App Security Data Connector to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

- `id` - The ID of the Microsoft Cloud App Security Data Connector.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

- `create` - (Defaults to 30 minutes) Used when creating the Microsoft Cloud App Security Data Connector.
- `read` - (Defaults to 5 minutes) Used when retrieving the Microsoft Cloud App Security Data Connector.
- `update` - (Defaults to 30 minutes) Used when updating the Microsoft Cloud App Security Data Connector.
- `delete` - (Defaults to 30 minutes) Used when deleting the Microsoft Cloud App Security Data Connector.

## Import

Microsoft Cloud App Security Data Connectors can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_sentinel_data_connector_microsoft_cloud_app_security.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.OperationalInsights/workspaces/workspace1/providers/Microsoft.SecurityInsights/dataConnectors/dc1
```
