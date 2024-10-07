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

* `data_collection_endpoint_id` - The resource ID of the Data Collection Endpoint that this rule can be used with.

* `data_flow` - One or more `data_flow` blocks as defined below.

* `destinations` - A `destinations` block as defined below.

* `location` - The Azure Region where the Data Collection Rule should exist. Changing this forces a new Data Collection Rule to be created.

* `data_sources` - A `data_sources` block as defined below. This property is optional and can be omitted if the rule is meant to be used via direct calls to the provisioned endpoint.

* `description` - The description of the Data Collection Rule.

* `identity` - An `identity` block as defined below.

* `kind` - The kind of the Data Collection Rule. Possible values are `Linux`, `Windows`,and `AgentDirectToStore`. A rule of kind `Linux` does not allow for `windows_event_log` data sources. And a rule of kind `Windows` does not allow for `syslog` data sources. If kind is not specified, all kinds of data sources are allowed.

* `stream_declaration` - A `stream_declaration` block as defined below.
 
* `tags` - A mapping of tags which should be assigned to the Data Collection Rule.

---

A `azure_monitor_metrics` block exports the following:

* `name` - The name which should be used for this destination. This name should be unique across all destinations regardless of type within the Data Collection Rule.

---

A `column` block exports the following:

* `name` - The name of the column.

* `type` - The type of the column data. Possible values are `string`, `int`, `long`, `real`, `boolean`, `datetime`,and `dynamic`.

---

A `data_import` block exports the following:

* `event_hub_data_source` - An `event_hub_data_source` block as defined below.

---

A `data_flow` block exports the following:

* `destinations` - Specifies a list of destination names. A `azure_monitor_metrics` data source only allows for stream of kind `Microsoft-InsightsMetrics`.

* `streams` - Specifies a list of streams. Possible values include but not limited to `Microsoft-Event`, `Microsoft-InsightsMetrics`, `Microsoft-Perf`, `Microsoft-Syslog`,and `Microsoft-WindowsEvent`.

* `built_in_transform` - The built-in transform to transform stream data.

* `output_stream` - The output stream of the transform. Only required if the data flow changes data to a different stream.

* `transform_kql` - The KQL query to transform stream data.

---

A `data_sources` block exports the following:

* `data_import` - A `data_import` block as defined above.

* `extension` - One or more `extension` blocks as defined below.

* `iis_log` - One or more `iis_log` blocks as defined below.

* `log_file` - One or more `log_file` blocks as defined below.

* `performance_counter` - One or more `performance_counter` blocks as defined below.

* `platform_telemetry` - One or more `platform_telemetry` blocks as defined below.

* `prometheus_forwarder` - One or more `prometheus_forwarder` blocks as defined below.

* `syslog` - One or more `syslog` blocks as defined below.

* `windows_event_log` - One or more `windows_event_log` blocks as defined below.

* `windows_firewall_log` - One or more `windows_firewall_log` blocks as defined below.

---

A `destinations` block exports the following:

* `azure_monitor_metrics` - A `azure_monitor_metrics` block as defined above.

* `event_hub` - One or more `event_hub` blocks as defined below.

* `event_hub_direct` - One or more `event_hub_direct` blocks as defined below.

* `log_analytics` - One or more `log_analytics` blocks as defined below.

* `monitor_account` - One or more `monitor_account` blocks as defined below.

* `storage_blob` - One or more `storage_blob` blocks as defined below.

* `storage_blob_direct` - One or more `storage_blob_direct` blocks as defined below.

* `storage_table_direct` - One or more `storage_table_direct` blocks as defined below.

---

An `event_hub_data_source` block exports the following:

* `name` - The name which should be used for this data source. This name should be unique across all data sources regardless of type within the Data Collection Rule.

* `stream` - The stream to collect from Event Hub. Possible value should be a custom stream name.

* `consumer_group` - The Event Hub consumer group name.

---

An `event_hub` block exports the following:

* `event_hub_id` - The resource ID of the Event Hub.

* `name` - The name which should be used for this destination. This name should be unique across all destinations regardless of type within the Data Collection Rule.

---

An `event_hub_direct` block exports the following:

* `event_hub_id` - The resource ID of the Event Hub.

* `name` - The name which should be used for this destination. This name should be unique across all destinations regardless of type within the Data Collection Rule.

---

A `extension` block exports the following:

* `extension_name` - The name of the VM extension.

* `name` - The name which should be used for this data source. This name should be unique across all data sources regardless of type within the Data Collection Rule.

* `streams` - Specifies a list of streams that this data source will be sent to. A stream indicates what schema will be used for this data and usually what table in Log Analytics the data will be sent to. Possible values include but not limited to `Microsoft-Event`, `Microsoft-InsightsMetrics`, `Microsoft-Perf`, `Microsoft-Syslog`,and `Microsoft-WindowsEvent`.

* `extension_json` - A JSON String which specifies the extension setting.

* `input_data_sources` - Specifies a list of data sources this extension needs data from. An item should be a name of a supported data source which produces only one stream. Supported data sources type: `performance_counter`, `windows_event_log`,and `syslog`.

---

An `iis_log` block exports the following:

* `name` - The name which should be used for this data source. This name should be unique across all data sources regardless of type within the Data Collection Rule.

* `streams` - Specifies a list of streams that this data source will be sent to. A stream indicates what schema will be used for this data and usually what table in Log Analytics the data will be sent to. Possible value is `Microsoft-W3CIISLog`.

* `log_directories` - Specifies a list of absolute paths where the log files are located.

---

An `identity` block exports the following:

* `type` - cSpecifies the type of Managed Service Identity that should be configured on this Data Collection Rule. Possible values are `SystemAssigned` and `UserAssigned`.

* `identity_ids` - A list of User Assigned Managed Identity IDs to be assigned to this Data Collection Rule. Currently, up to 1 identity is supported.

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

---

A `label_include_filter` block exports the following:

* `label` - The label of the filter. This label should be unique across all `label_include_fileter` block. Possible value is `microsoft_metrics_include_label`.

* `value` - The value of the filter.

---

An `log_file` block exports the following:

* `name` - The name which should be used for this data source. This name should be unique across all data sources regardless of type within the Data Collection Rule.

* `streams` - Specifies a list of streams that this data source will be sent to. A stream indicates what schema will be used for this data and usually what table in Log Analytics the data will be sent to. Possible value should be custom stream names.

* `file_patterns` - Specifies a list of file patterns where the log files are located. For example, `C:\\JavaLogs\\*.log`.

* `format` - The data format of the log files. possible value is `text`.

* `settings` - A `settings` block as defined below.

---

A `monitor_account` block exports the following:

* `monitor_account_id` - The resource ID of the Monitor Account.

* `name` - The name which should be used for this destination. This name should be unique across all destinations regardless of type within the Data Collection Rule.

---
A `log_analytics` block exports the following:

* `name` - The name which should be used for this destination. This name should be unique across all destinations regardless of type within the Data Collection Rule.

* `workspace_resource_id` - The ID of a Log Analytic Workspace resource.

---

A `performance_counter` block exports the following:

* `counter_specifiers` - Specifies a list of specifier names of the performance counters you want to collect. Use a wildcard `*` to collect counters for all instances. To get a list of performance counters on Windows, run the command `typeperf`.

* `name` - The name which should be used for this data source. This name should be unique across all data sources regardless of type within the Data Collection Rule.

* `sampling_frequency_in_seconds` - The number of seconds between consecutive counter measurements (samples). The value should be integer between `1` and `1800` inclusive.

* `streams` - Specifies a list of streams that this data source will be sent to. A stream indicates what schema will be used for this data and usually what table in Log Analytics the data will be sent to. Possible values include but not limited to `Microsoft-InsightsMetrics`,and `Microsoft-Perf`.

---

A `platform_telemetry` block exports the following:

* `name` - The name which should be used for this data source. This name should be unique across all data sources regardless of type within the Data Collection Rule.

* `streams` - Specifies a list of streams that this data source will be sent to. A stream indicates what schema will be used for this data and usually what table in Log Analytics the data will be sent to. Possible values include but not limited to `Microsoft.Cache/redis:Metrics-Group-All`.

---

A `prometheus_forwarder` block exports the following:

* `name` - The name which should be used for this data source. This name should be unique across all data sources regardless of type within the Data Collection Rule.

* `streams` - Specifies a list of streams that this data source will be sent to. A stream indicates what schema will be used for this data and usually what table in Log Analytics the data will be sent to. Possible value is `Microsoft-PrometheusMetrics`.

* `label_include_filter` - One or more `label_include_filter` blocks as defined above.

---

A `settings` block within the `log_file` block exports the following:

* `text` - A `text` block as defined below.

---

A `storage_blob` block exports the following:

* `container_name` - The Storage Container name.

* `name` - The name which should be used for this destination. This name should be unique across all destinations regardless of type within the Data Collection Rule.

* `storage_account_id` - The resource ID of the Storage Account.

---

A `storage_blob_direct` block exports the following:

* `container_name` - The Storage Container name.

* `name` - The name which should be used for this destination. This name should be unique across all destinations regardless of type within the Data Collection Rule.

* `storage_account_id` - The resource ID of the Storage Account.

---

A `storage_table_direct` block exports the following:

* `table_name` - The Storage Table name.

* `name` - The name which should be used for this destination. This name should be unique across all destinations regardless of type within the Data Collection Rule.

* `storage_account_id` - The resource ID of the Storage Account.

---

A `stream_declaration` block exports the following:

* `stream_name` - The name of the custom stream. This name should be unique across all `stream_declaration` blocks.

* `column` - One or more `column` blocks as defined above.

---

A `syslog` block exports the following:

* `facility_names` - Specifies a list of facility names. Use a wildcard `*` to collect logs for all facility names. Possible values are `auth`, `authpriv`, `cron`, `daemon`, `kern`, `lpr`, `mail`, `mark`, `news`, `syslog`, `user`, `uucp`, `local0`, `local1`, `local2`, `local3`, `local4`, `local5`, `local6`, `local7`,and `*`.

* `log_levels` - Specifies a list of log levels. Use a wildcard `*` to collect logs for all log levels. Possible values are `Debug`,  `Info`, `Notice`, `Warning`, `Error`, `Critical`, `Alert`, `Emergency`,and `*`.

* `name` - The name which should be used for this data source. This name should be unique across all data sources regardless of type within the Data Collection Rule.

* `streams` - Specifies a list of streams that this data source will be sent to. A stream indicates what schema will be used for this data and usually what table in Log Analytics the data will be sent to. Possible values include but not limited to `Microsoft-Syslog`,and `Microsoft-CiscoAsa`, and `Microsoft-CommonSecurityLog`.

---

A `text` block within the `log_file.settings` block exports the following:

* `record_start_timestamp_format` - The timestamp format of the text log files. Possible values are `ISO 8601`, `YYYY-MM-DD HH:MM:SS`, `M/D/YYYY HH:MM:SS AM/PM`, `Mon DD, YYYY HH:MM:SS`, `yyMMdd HH:mm:ss`, `ddMMyy HH:mm:ss`, `MMM d hh:mm:ss`, `dd/MMM/yyyy:HH:mm:ss zzz`,and `yyyy-MM-ddTHH:mm:ssK`.

---

A `windows_event_log` block exports the following:

* `name` - The name which should be used for this data source. This name should be unique across all data sources regardless of type within the Data Collection Rule.

* `streams` - Specifies a list of streams that this data source will be sent to. A stream indicates what schema will be used for this data and usually what table in Log Analytics the data will be sent to. Possible values include but not limited to `Microsoft-Event`,and `Microsoft-WindowsEvent`, `Microsoft-RomeDetectionEvent`, and `Microsoft-SecurityEvent`.

* `x_path_queries` - Specifies a list of Windows Event Log queries in XPath expression.

---

A `windows_firewall_log` block exports the following:

* `name` - The name which should be used for this data source. This name should be unique across all data sources regardless of type within the Data Collection Rule.

* `streams` - Specifies a list of streams that this data source will be sent to. A stream indicates what schema will be used for this data and usually what table in Log Analytics the data will be sent to.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Data Collection Rule.
