---
subcategory: "Sentinel"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_sentinel_data_connector_office_365"
description: |-
  Manages a Office 365 Data Connector.
---

# azurerm_sentinel_data_connector_office_365

Manages a Office 365 Data Connector.

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

resource "azurerm_sentinel_data_connector_office_365" "example" {
  name                       = "example"
  log_analytics_workspace_id = azurerm_log_analytics_workspace.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `log_analytics_workspace_id` - (Required) The ID of the Log Analytics Workspace that this Office 365 Data Connector resides in. Changing this forces a new Office 365 Data Connector to be created.

* `name` - (Required) The name which should be used for this Office 365 Data Connector. Changing this forces a new Office 365 Data Connector to be created.

---

* `exchange_enabled` - (Optional) Should the Exchange data connector be enabled? Defaults to `true`.

* `sharepoint_enabled` - (Optional) Should the SharePoint data connector be enabled? Defaults to `true`.

* `teams_enabled` - (Optional) Should the Microsoft Teams data connector be enabled? Defaults to `true`.

-> **NOTE**: At least one of `exchange_enabled`, `sharedpoint_enabled` and `teams_enabled` has to be specified.

* `tenant_id` - (Optional) The ID of the Tenant that this Office 365 Data Connector connects to. Changing this forces a new Office 365 Data Connector to be created.

-> **NOTE**: Terraform will use the Tenant ID for the current Subscription if this is unspecified.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Office 365 Data Connector.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Office 365 Data Connector.
* `read` - (Defaults to 5 minutes) Used when retrieving the Office 365 Data Connector.
* `update` - (Defaults to 30 minutes) Used when updating the Office 365 Data Connector.
* `delete` - (Defaults to 30 minutes) Used when deleting the Office 365 Data Connector.

## Import

Office 365 Data Connectors can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_sentinel_data_connector_office_365.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.OperationalInsights/workspaces/workspace1/providers/Microsoft.SecurityInsights/dataConnectors/dc1
```
