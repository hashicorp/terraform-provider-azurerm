---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_monitor_diagnostic_setting"
sidebar_current: "docs-azurerm-resource-monitor-diagnostic-setting"
description: |-
  Manages a Diagnostic Setting for an existing Resource.

---

# azurerm_monitor_diagnostic_setting

Manages a Diagnostic Setting for an existing Resource.

## Example Usage

```

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Diagnostic Setting. Changing this forces a new resource to be created.

* `resource_id` - (Required) The ID of an existing Resource on which to configure Diagnostic Settings. Changing this forces a new resource to be created.

* `event_hub_name` - (Optional) Specifies the name of the Event Hub where Diagnostics Data should be sent. Changing this forces a new resource to be created.

-> **NOTE:** If this isn't specified then the default Event Hub will be used.

* `event_hub_authorization_rule_id` - (Optional) Specifies the ID of an Event Hub Namespace Authorization Rule used to send Diagnostics Data. Changing this forces a new resource to be created.

-> **NOTE:** One of `event_hub_authorization_rule_id`, `log_analytics_workspace_id` and `storage_account_id` must be specified.

* `log` - (Optional) One or more `log` blocks as defined below.

-> **NOTE:** At least one `log` or `metric` block must be specified.

* `log_analytics_workspace_id` - (Optional) Specifies the ID of a Log Analytics Workspace where Diagnostics Data should be sent. Changing this forces a new resource to be created.

-> **NOTE:** One of `event_hub_authorization_rule_id`, `log_analytics_workspace_id` and `storage_account_id` must be specified.

* `metric` - (Optional) One or more `metric` blocks as defined below.

-> **NOTE:** At least one `log` or `metric` block must be specified.

* `storage_account_id` - (Optional) With this parameter you can specify a storage account which should be used to send the logs to. Parameter must be a valid Azure Resource ID. Changing this forces a new resource to be created.

-> **NOTE:** One of `event_hub_authorization_rule_id`, `log_analytics_workspace_id` and `storage_account_id` must be specified.

---

A `log` block supports the following:

* `category` - (Required) The name of a Diagnostic Log Category for this Resource.

-> **NOTE:** The Log Categories available vary depending on the Resource being used. You may wish to use [the `azurerm_monitor_diagnostic_categories` Data Source](../d/monitor_diagnostic_categories.html) to identify which categories are available for a given Resource.

* `retention_policy` - (Required) A `retention_policy` block as defined below.

* `enabled` - (Optional) Is this Diagnostic Log enabled? Defaults to `true`.

---

A `metric` block supports the following:

* `category` - (Required) The name of a Diagnostic Metric Category for this Resource.

-> **NOTE:** The Metric Categories available vary depending on the Resource being used. You may wish to use [the `azurerm_monitor_diagnostic_categories` Data Source](../d/monitor_diagnostic_categories.html) to identify which categories are available for a given Resource.

* `retention_policy` - (Required) A `retention_policy` block as defined below.

* `enabled` - (Optional) Is this Diagnostic Metric enabled? Defaults to `true`.

---

A `retention_policy` block supports the following:

* `enabled` - (Required) Is this Retention Policy enabled?

* `days` - (Optional) The number of days for which this Retention Policy should apply.

-> **NOTE:** Setting this to `0` will retain the events indefinitely.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Diagnostic Setting.

## Import

Diagnostic Settings can be imported using the `resource id`, e.g.

```
terraform import azurerm_monitor_diagnostics.test /subscriptions/XXX/resourcegroups/resource_group/providers/microsoft.keyvault/vaults/vault|logMonitoring
```

-> **NOTE:** This is a Terraform specific Resource ID which uses the format `{resourceId}|{diagnosticSettingName}`