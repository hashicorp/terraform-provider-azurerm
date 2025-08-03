---
subcategory: "Sentinel"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_sentinel_data_connector_office_power_bi"
description: |-
  Manages an Office Power BI Data Connector.
---

# azurerm_sentinel_data_connector_office_power_bi

Manages an Office Power BI Data Connector.

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

resource "azurerm_sentinel_data_connector_office_power_bi" "example" {
  name                       = "example"
  log_analytics_workspace_id = azurerm_sentinel_log_analytics_workspace_onboarding.example.workspace_id
}
```

## Arguments Reference

The following arguments are supported:

* `log_analytics_workspace_id` - (Required) The ID of the Log Analytics Workspace that this Office Power BI Data Connector resides in. Changing this forces a new Office Power BI Data Connector to be created.

* `name` - (Required) The name which should be used for this Office Power BI Data Connector. Changing this forces a new Office Power BI Data Connector to be created.

---

* `tenant_id` - (Optional) The ID of the tenant that this Office Power BI Data Connector connects to. Changing this forces a new Office Power BI Data Connector to be created.

-> **Note:** Currently, only the same tenant as the running account is allowed. Cross-tenant scenario is not supported yet.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Office Power BI Data Connector.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Office Power BI Data Connector.
* `read` - (Defaults to 5 minutes) Used when retrieving the Office Power BI Data Connector.
* `delete` - (Defaults to 30 minutes) Used when deleting the Office Power BI Data Connector.

## Import

Office Power BI Data Connectors can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_sentinel_data_connector_office_power_bi.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.OperationalInsights/workspaces/workspace1/providers/Microsoft.SecurityInsights/dataConnectors/dc1
```
