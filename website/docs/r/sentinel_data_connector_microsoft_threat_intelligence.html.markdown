---
subcategory: "Sentinel"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_sentinel_data_connector_microsoft_threat_intelligence"
description: |-
  Manages a Microsoft Threat Intelligence Data Connector.
---

# azurerm_sentinel_data_connector_microsoft_threat_intelligence

Manages a Microsoft Threat Intelligence Data Connector.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "east us"
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

resource "azurerm_sentinel_data_connector_microsoft_threat_intelligence" "example" {
  name                                         = "example-dc-msti"
  log_analytics_workspace_id                   = azurerm_sentinel_log_analytics_workspace_onboarding.example.workspace_id
  microsoft_emerging_threat_feed_lookback_date = "1970-01-01T00:00:00Z"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Microsoft Threat Intelligence Data Connector. Changing this forces a new Microsoft Threat Intelligence Data Connector to be created.

* `log_analytics_workspace_id` - (Required) The ID of the Log Analytics Workspace. Changing this forces a new Data Connector to be created.

* `microsoft_emerging_threat_feed_lookback_date` - (Required) The lookback date for the Microsoft Emerging Threat Feed in RFC3339. Changing this forces a new Data Connector to be created.

---

* `tenant_id` - (Optional) The ID of the tenant that this Microsoft Threat Intelligence Data Connector connects to. Changing this forces a new Microsoft Threat Intelligence Data Connector to be created.

-> **Note:** Currently, only the same tenant as the running account is allowed. Cross-tenant scenario is not supported yet.


## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the sentinel.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the sentinel.
* `read` - (Defaults to 5 minutes) Used when retrieving the sentinel.
* `delete` - (Defaults to 30 minutes) Used when deleting the sentinel.

## Import

sentinels can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_sentinel_data_connector_microsoft_threat_intelligence.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.OperationalInsights/workspaces/workspace1/providers/Microsoft.SecurityInsights/dataConnectors/dc1
```
