---
subcategory: "Monitor"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_monitor_data_collection_rule"
description: |-
  Get information about the specified Data Collection Rule.

---

# Data Source: azurerm_monitor_data_collection_rule

Use this data source to access information about an existing Data Collection Rule.

## Example Usage

```hcl
data "azurerm_monitor_data_collection_rule" "example" {
  name                = "example-rule"
  resource_group_name = azurerm_resource_group.example.name
}

output "rule_id" {
  value = data.azurerm_monitor_data_collection_rule.example.id
}
```

## Argument Reference

* `name` - Specifies the name of the Data Collection Rule.

* `resource_group_name` - Specifies the name of the resource group the Data Collection Rule is located in.

## Attributes Reference

* `id` - The ID of the Resource.

* `data_flow` - One or more `data_flow` blocks as defined below.

* `destinations` - A `destinations` block as defined below.

* `location` - The Azure Region where the Data Collection Rule should exist. Changing this forces a new Data Collection Rule to be created.

* `data_sources` - A `data_sources` block as defined below. This property is optional and can be omitted if the rule is meant to be used via direct calls to the provisioned endpoint.

* `description` - The description of the Data Collection Rule.

* `kind` - The kind of the Data Collection Rule. Possible values are `Linux` and `Windows`. A rule of kind `Linux` does not allow for `windows_event_log` data sources. And a rule of kind `Windows` does not allow for `syslog` data sources. If kind is not specified, all kinds of data sources are allowed.

* `tags` - A mapping of tags which should be assigned to the Data Collection Rule.

---

A `azure_monitor_metrics` block supports the following:

* `name` - The name which should be used for this destination. This name should be unique across all destinations regardless of type within the Data Collection Rule.

---

A `data_flow` block supports the following:

* `destinations` - Specifies a list of destination names. A `azure_monitor_metrics` data source only allows for stream of kind `Microsoft-InsightsMetrics`.

* `streams` - Specifies a list of streams. Possible values are `Microsoft-Event`, `Microsoft-InsightsMetrics`, `Microsoft-Perf`, `Microsoft-Syslog`,and `Microsoft-WindowsEvent`.

---

A `data_sources` block supports the following:

* `extension` - One or more `extension` blocks as defined below.

* `performance_counter` - One or more `performance_counter` blocks as defined below.

* `syslog` - One or more `syslog` blocks as defined below.

* `windows_event_log` - One or more `windows_event_log` blocks as defined below.

---

A `destinations` block supports the following:

* `azure_monitor_metrics` - A `azure_monitor_metrics` block as defined above.

* `log_analytics` - One or more `log_analytics` blocks as defined below.

---

A `extension` block supports the following:

* `extension_name` - The name of the VM extension.

* `name` - The name which should be used for this data source. This name should be unique across all data sources regardless of type within the Data Collection Rule.

* `streams` - Specifies a list of streams that this data source will be sent to. A stream indicates what schema will be used for this data and usually what table in Log Analytics the data will be sent to. Possible values are `Microsoft-Event`, `Microsoft-InsightsMetrics`, `Microsoft-Perf`, `Microsoft-Syslog`,and `Microsoft-WindowsEvent`.

* `extension_json` - A JSON String which specifies the extension setting.

* `input_data_sources` - Specifies a list of data sources this extension needs data from. An item should be a name of a supported data source which produces only one stream. Supported data sources type: `performance_counter`, `windows_event_log`,and `syslog`.

---

A `log_analytics` block supports the following:

* `name` - The name which should be used for this destination. This name should be unique across all destinations regardless of type within the Data Collection Rule.

* `workspace_resource_id` - The ID of a Log Analytic Workspace resource.

---

A `performance_counter` block supports the following:

* `counter_specifiers` - Specifies a list of specifier names of the performance counters you want to collect. Use a wildcard `*` to collect counters for all instances. To get a list of performance counters on Windows, run the command `typeperf`.

* `name` - The name which should be used for this data source. This name should be unique across all data sources regardless of type within the Data Collection Rule.

* `sampling_frequency_in_seconds` - The number of seconds between consecutive counter measurements (samples). The value should be integer between `1` and `300` inclusive.

* `streams` - Specifies a list of streams that this data source will be sent to. A stream indicates what schema will be used for this data and usually what table in Log Analytics the data will be sent to. Possible values are `Microsoft-InsightsMetrics`,and `Microsoft-Perf`.

---

A `syslog` block supports the following:

* `facility_names` - Specifies a list of facility names. Use a wildcard `*` to collect logs for all facility names. Possible values are `auth`, `authpriv`, `cron`, `daemon`, `kern`, `lpr`, `mail`, `mark`, `news`, `syslog`, `user`, `uucp`, `local0`, `local1`, `local2`, `local3`, `local4`, `local5`, `local6`, `local7`,and `*`.

* `log_levels` - Specifies a list of log levels. Use a wildcard `*` to collect logs for all log levels. Possible values are `Debug`,  `Info`, `Notice`, `Warning`, `Error`, `Critical`, `Alert`, `Emergency`,and `*`.

* `name` - The name which should be used for this data source. This name should be unique across all data sources regardless of type within the Data Collection Rule.

-> **Note:** Syslog data source has only one possible streams value which is `Microsoft-Syslog`.

---

A `windows_event_log` block supports the following:

* `name` - The name which should be used for this data source. This name should be unique across all data sources regardless of type within the Data Collection Rule.

* `streams` - Specifies a list of streams that this data source will be sent to. A stream indicates what schema will be used for this data and usually what table in Log Analytics the data will be sent to. Possible values are `Microsoft-Event`,and `Microsoft-WindowsEvent`.

* `x_path_queries` - Specifies a list of Windows Event Log queries in XPath expression.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Data Collection Rule.
