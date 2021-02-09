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

data "azurerm_storage_account" "example" {
  name                = "examplestoracc"
  resource_group_name = azurerm_resource_group.example.name
}

data "azurerm_key_vault" "example" {
  name                = "example-vault"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_monitor_diagnostic_setting" "example" {
  name               = "example"
  target_resource_id = data.azurerm_key_vault.example.id
  storage_account_id = data.azurerm_storage_account.example.id

  log {
    category = "AuditEvent"
    enabled  = false

    retention_policy {
      enabled = false
    }
  }

  metric {
    category = "AllMetrics"

    retention_policy {
      enabled = false
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Diagnostic Setting. Changing this forces a new resource to be created.

* `target_resource_id` - (Required) The ID of an existing Resource on which to configure Diagnostic Settings. Changing this forces a new resource to be created.

* `eventhub_name` - (Optional) Specifies the name of the Event Hub where Diagnostics Data should be sent. Changing this forces a new resource to be created.

-> **NOTE:** If this isn't specified then the default Event Hub will be used.

* `eventhub_authorization_rule_id` - (Optional) Specifies the ID of an Event Hub Namespace Authorization Rule used to send Diagnostics Data. Changing this forces a new resource to be created.

-> **NOTE:** This can be sourced from [the `azurerm_eventhub_namespace_authorization_rule` resource](eventhub_namespace_authorization_rule.html) and is different from [a `azurerm_eventhub_authorization_rule` resource](eventhub_authorization_rule.html).

-> **NOTE:** One of `eventhub_authorization_rule_id`, `log_analytics_workspace_id` and `storage_account_id` must be specified.

* `log` - (Optional) One or more `log` blocks as defined below.

-> **NOTE:** At least one `log` or `metric` block must be specified.

* `log_analytics_workspace_id` - (Optional) Specifies the ID of a Log Analytics Workspace where Diagnostics Data should be sent.

-> **NOTE:** One of `eventhub_authorization_rule_id`, `log_analytics_workspace_id` and `storage_account_id` must be specified.

* `metric` - (Optional) One or more `metric` blocks as defined below.

-> **NOTE:** At least one `log` or `metric` block must be specified.

* `storage_account_id` - (Optional) The ID of the Storage Account where logs should be sent. Changing this forces a new resource to be created.

-> **NOTE:** One of `eventhub_authorization_rule_id`, `log_analytics_workspace_id` and `storage_account_id` must be specified.

* `log_analytics_destination_type` - (Optional) When set to 'Dedicated' logs sent to a Log Analytics workspace will go into resource specific tables, instead of the legacy AzureDiagnostics table.

-> **NOTE:** This setting will only have an effect if a `log_analytics_workspace_id` is provided, and the resource is available for resource-specific logs.  As of July 2019, this only includes Azure Data Factory. Please [see the documentation](https://docs.microsoft.com/en-us/azure/azure-monitor/platform/diagnostic-logs-stream-log-store#azure-diagnostics-vs-resource-specific) for more information.

---

A `log` block supports the following:

* `category` - (Required) The name of a Diagnostic Log Category for this Resource.

-> **NOTE:** The Log Categories available vary depending on the Resource being used. You may wish to use [the `azurerm_monitor_diagnostic_categories` Data Source](../d/monitor_diagnostic_categories.html) or [list of service specific schemas](https://docs.microsoft.com/en-us/azure/azure-monitor/platform/resource-logs-schema#service-specific-schemas) to identify which categories are available for a given Resource.

* `retention_policy` - (Optional) A `retention_policy` block as defined below.

* `enabled` - (Optional) Is this Diagnostic Log enabled? Defaults to `true`.

---

A `metric` block supports the following:

* `category` - (Required) The name of a Diagnostic Metric Category for this Resource.

-> **NOTE:** The Metric Categories available vary depending on the Resource being used. You may wish to use [the `azurerm_monitor_diagnostic_categories` Data Source](../d/monitor_diagnostic_categories.html) to identify which categories are available for a given Resource.

* `retention_policy` - (Optional) A `retention_policy` block as defined below.

* `enabled` - (Optional) Is this Diagnostic Metric enabled? Defaults to `true`.

---

A `retention_policy` block supports the following:

* `enabled` - (Required) Is this Retention Policy enabled?

* `days` - (Optional) The number of days for which this Retention Policy should apply.

-> **NOTE:** Setting this to `0` will retain the events indefinitely.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Diagnostic Setting.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Diagnostics Setting.
* `update` - (Defaults to 30 minutes) Used when updating the Diagnostics Setting.
* `read` - (Defaults to 5 minutes) Used when retrieving the Diagnostics Setting.
* `delete` - (Defaults to 60 minutes) Used when deleting the Diagnostics Setting.

## Import

Diagnostic Settings can be imported using the `resource id`, e.g.

```
terraform import azurerm_monitor_diagnostic_setting.example "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.KeyVault/vaults/vault1|logMonitoring1"
```

-> **NOTE:** This is a Terraform specific Resource ID which uses the format `{resourceId}|{diagnosticSettingName}`
