---
subcategory: "Sentinel"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_sentinel_data_connector_iot"
description: |-
  Manages an Iot Data Connector.
---

# azurerm_sentinel_data_connector_iot

Manages an Iot Data Connector.

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

resource "azurerm_sentinel_log_analytics_workspace_onboarding" "example" {
  workspace_id = azurerm_log_analytics_workspace.example.id
}

resource "azurerm_sentinel_data_connector_iot" "example" {
  name                       = "example"
  log_analytics_workspace_id = azurerm_sentinel_log_analytics_workspace_onboarding.example.workspace_id
}
```

## Arguments Reference

The following arguments are supported:

* `log_analytics_workspace_id` - (Required) The ID of the Log Analytics Workspace that this Iot Data Connector resides in. Changing this forces a new Iot Data Connector to be created.

* `name` - (Required) The name which should be used for this Iot Data Connector. Changing this forces a new Iot Data Connector to be created.

---

* `subscription_id` - (Optional) The ID of the subscription that this Iot Data Connector connects to. Changing this forces a new Iot Data Connector to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Iot Data Connector.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Iot Data Connector.
* `read` - (Defaults to 5 minutes) Used when retrieving the Iot Data Connector.
* `delete` - (Defaults to 30 minutes) Used when deleting the Iot Data Connector.

## Import

Iot Data Connectors can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_sentinel_data_connector_iot.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.OperationalInsights/workspaces/workspace1/providers/Microsoft.SecurityInsights/dataConnectors/dc1
```
