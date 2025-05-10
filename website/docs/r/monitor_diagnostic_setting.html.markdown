---
subcategory: "Monitor"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_monitor_diagnostic_setting"
description: |-
  Manages a Diagnostic Setting for an existing Resource.

---

# azurerm_monitor_diagnostic_setting

Manages a Diagnostic Setting for an existing Resource.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "storageaccountname"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "example" {
  name                       = "examplekeyvault"
  location                   = azurerm_resource_group.example.location
  resource_group_name        = azurerm_resource_group.example.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  soft_delete_retention_days = 7
  purge_protection_enabled   = false
  sku_name                   = "standard"
}

resource "azurerm_monitor_diagnostic_setting" "example" {
  name               = "example"
  target_resource_id = azurerm_key_vault.example.id
  storage_account_id = azurerm_storage_account.example.id

  enabled_log {
    category = "AuditEvent"
  }

  metric {
    category = "AllMetrics"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Diagnostic Setting. Changing this forces a new resource to be created.

-> **Note:** If the name is set to 'service' it will not be possible to fully delete the diagnostic setting. This is due to legacy API support.

* `target_resource_id` - (Required) The ID of an existing Resource on which to configure Diagnostic Settings. Changing this forces a new resource to be created.

* `eventhub_name` - (Optional) Specifies the name of the Event Hub where Diagnostics Data should be sent.

-> **Note:** If this isn't specified then the default Event Hub will be used.

* `eventhub_authorization_rule_id` - (Optional) Specifies the ID of an Event Hub Namespace Authorization Rule used to send Diagnostics Data. 

-> **Note:** This can be sourced from [the `azurerm_eventhub_namespace_authorization_rule` resource](eventhub_namespace_authorization_rule.html) and is different from [a `azurerm_eventhub_authorization_rule` resource](eventhub_authorization_rule.html).

-> **Note:** At least one of `eventhub_authorization_rule_id`, `log_analytics_workspace_id`, `partner_solution_id` and `storage_account_id` must be specified.

* `enabled_log` - (Optional) One or more `enabled_log` blocks as defined below.

-> **Note:** At least one `enabled_log` or `metric` block must be specified. At least one type of Log or Metric must be enabled.

* `log_analytics_workspace_id` - (Optional) Specifies the ID of a Log Analytics Workspace where Diagnostics Data should be sent.

-> **Note:** At least one of `eventhub_authorization_rule_id`, `log_analytics_workspace_id`, `partner_solution_id` and `storage_account_id` must be specified.

* `metric` - (Optional) One or more `metric` blocks as defined below.

-> **Note:** At least one `enabled_log` or `metric` block must be specified.

* `storage_account_id` - (Optional) The ID of the Storage Account where logs should be sent. 

-> **Note:** At least one of `eventhub_authorization_rule_id`, `log_analytics_workspace_id`, `partner_solution_id` and `storage_account_id` must be specified.

* `log_analytics_destination_type` - (Optional) Possible values are `AzureDiagnostics` and `Dedicated`. When set to `Dedicated`, logs sent to a Log Analytics workspace will go into resource specific tables, instead of the legacy `AzureDiagnostics` table.

-> **Note:** This setting will only have an effect if a `log_analytics_workspace_id` is provided. For some target resource type (e.g., Key Vault), this field is unconfigurable. Please see [resource types](https://learn.microsoft.com/en-us/azure/azure-monitor/reference/tables/azurediagnostics#resource-types) for services that use each method. Please [see the documentation](https://docs.microsoft.com/azure/azure-monitor/platform/diagnostic-logs-stream-log-store#azure-diagnostics-vs-resource-specific) for details on the differences between destination types.

* `partner_solution_id` - (Optional) The ID of the market partner solution where Diagnostics Data should be sent. For potential partner integrations, [click to learn more about partner integration](https://learn.microsoft.com/en-us/azure/partner-solutions/overview).

-> **Note:** At least one of `eventhub_authorization_rule_id`, `log_analytics_workspace_id`, `partner_solution_id` and `storage_account_id` must be specified.

---

An `enabled_log` block supports the following:

* `category` - (Optional) The name of a Diagnostic Log Category for this Resource.

-> **Note:** The Log Categories available vary depending on the Resource being used. You may wish to use [the `azurerm_monitor_diagnostic_categories` Data Source](../d/monitor_diagnostic_categories.html) or [list of service specific schemas](https://docs.microsoft.com/azure/azure-monitor/platform/resource-logs-schema#service-specific-schemas) to identify which categories are available for a given Resource.

* `category_group` - (Optional) The name of a Diagnostic Log Category Group for this Resource.

-> **Note:** Not all resources have category groups available.

-> **Note:** Exactly one of `category` or `category_group` must be specified.

---

A `metric` block supports the following:

* `category` - (Required) The name of a Diagnostic Metric Category for this Resource.

-> **Note:** The Metric Categories available vary depending on the Resource being used. You may wish to use [the `azurerm_monitor_diagnostic_categories` Data Source](../d/monitor_diagnostic_categories.html) to identify which categories are available for a given Resource.

* `enabled` - (Optional) Is this Diagnostic Metric enabled? Defaults to `true`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Diagnostic Setting.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Diagnostics Setting.
* `read` - (Defaults to 5 minutes) Used when retrieving the Diagnostics Setting.
* `update` - (Defaults to 30 minutes) Used when updating the Diagnostics Setting.
* `delete` - (Defaults to 1 hour) Used when deleting the Diagnostics Setting.

## Import

Diagnostic Settings can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_monitor_diagnostic_setting.example "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.KeyVault/vaults/vault1|logMonitoring1"
```

-> **Note:** This is a Terraform specific Resource ID which uses the format `{resourceId}|{diagnosticSettingName}`
