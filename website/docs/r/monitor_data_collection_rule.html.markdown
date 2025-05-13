---
subcategory: "Monitor"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_monitor_data_collection_rule"
description: |-
  Manages a Data Collection Rule.
---

# azurerm_monitor_data_collection_rule

Manages a Data Collection Rule.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_user_assigned_identity" "example" {
  name                = "example-uai"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_log_analytics_workspace" "example" {
  name                = "example-workspace"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_log_analytics_solution" "example" {
  solution_name         = "WindowsEventForwarding"
  location              = azurerm_resource_group.example.location
  resource_group_name   = azurerm_resource_group.example.name
  workspace_resource_id = azurerm_log_analytics_workspace.example.id
  workspace_name        = azurerm_log_analytics_workspace.example.name
  plan {
    publisher = "Microsoft"
    product   = "OMSGallery/WindowsEventForwarding"
  }
}

resource "azurerm_eventhub_namespace" "example" {
  name                = "exeventns"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "Standard"
  capacity            = 1
}

resource "azurerm_eventhub" "example" {
  name                = "exevent2"
  namespace_name      = azurerm_eventhub_namespace.example.name
  resource_group_name = azurerm_resource_group.example.name
  partition_count     = 2
  message_retention   = 1
}

resource "azurerm_storage_account" "example" {
  name                     = "examstorage"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "example" {
  name                  = "examplecontainer"
  storage_account_name  = azurerm_storage_account.example.name
  container_access_type = "private"
}

resource "azurerm_monitor_data_collection_endpoint" "example" {
  name                = "example-dcre"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  lifecycle {
    create_before_destroy = true
  }
}

resource "azurerm_monitor_data_collection_rule" "example" {
  name                        = "example-rule"
  resource_group_name         = azurerm_resource_group.example.name
  location                    = azurerm_resource_group.example.location
  data_collection_endpoint_id = azurerm_monitor_data_collection_endpoint.example.id

  destinations {
    log_analytics {
      workspace_resource_id = azurerm_log_analytics_workspace.example.id
      name                  = "example-destination-log"
    }

    event_hub {
      event_hub_id = azurerm_eventhub.example.id
      name         = "example-destination-eventhub"
    }

    storage_blob {
      storage_account_id = azurerm_storage_account.example.id
      container_name     = azurerm_storage_container.example.name
      name               = "example-destination-storage"
    }

    azure_monitor_metrics {
      name = "example-destination-metrics"
    }
  }

  data_flow {
    streams      = ["Microsoft-InsightsMetrics"]
    destinations = ["example-destination-metrics"]
  }

  data_flow {
    streams      = ["Microsoft-InsightsMetrics", "Microsoft-Syslog", "Microsoft-Perf"]
    destinations = ["example-destination-log"]
  }

  data_flow {
    streams       = ["Custom-MyTableRawData"]
    destinations  = ["example-destination-log"]
    output_stream = "Microsoft-Syslog"
    transform_kql = "source | project TimeGenerated = Time, Computer, Message = AdditionalContext"
  }

  data_sources {
    syslog {
      facility_names = ["*"]
      log_levels     = ["*"]
      name           = "example-datasource-syslog"
      streams        = ["Microsoft-Syslog"]
    }

    iis_log {
      streams         = ["Microsoft-W3CIISLog"]
      name            = "example-datasource-iis"
      log_directories = ["C:\\Logs\\W3SVC1"]
    }

    log_file {
      name          = "example-datasource-logfile"
      format        = "text"
      streams       = ["Custom-MyTableRawData"]
      file_patterns = ["C:\\JavaLogs\\*.log"]
      settings {
        text {
          record_start_timestamp_format = "ISO 8601"
        }
      }
    }

    performance_counter {
      streams                       = ["Microsoft-Perf", "Microsoft-InsightsMetrics"]
      sampling_frequency_in_seconds = 60
      counter_specifiers            = ["Processor(*)\\% Processor Time"]
      name                          = "example-datasource-perfcounter"
    }

    windows_event_log {
      streams        = ["Microsoft-WindowsEvent"]
      x_path_queries = ["*![System/Level=1]"]
      name           = "example-datasource-wineventlog"
    }

    extension {
      streams            = ["Microsoft-WindowsEvent"]
      input_data_sources = ["example-datasource-wineventlog"]
      extension_name     = "example-extension-name"
      extension_json = jsonencode({
        a = 1
        b = "hello"
      })
      name = "example-datasource-extension"
    }
  }

  stream_declaration {
    stream_name = "Custom-MyTableRawData"
    column {
      name = "Time"
      type = "datetime"
    }
    column {
      name = "Computer"
      type = "string"
    }
    column {
      name = "AdditionalContext"
      type = "string"
    }
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.example.id]
  }

  description = "data collection rule example"
  tags = {
    foo = "bar"
  }
  depends_on = [
    azurerm_log_analytics_solution.example
  ]
}
```

## Arguments Reference

The following arguments are supported:

* `data_flow` - (Required) One or more `data_flow` blocks as defined below.

* `destinations` - (Required) A `destinations` block as defined below.

* `location` - (Required) The Azure Region where the Data Collection Rule should exist. Changing this forces a new Data Collection Rule to be created.

* `name` - (Required) The name which should be used for this Data Collection Rule. Changing this forces a new Data Collection Rule to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Data Collection Rule should exist. Changing this forces a new Data Collection Rule to be created.

---

* `data_collection_endpoint_id` - (Optional) The resource ID of the Data Collection Endpoint that this rule can be used with.

* `data_sources` - (Optional) A `data_sources` block as defined below. This property is optional and can be omitted if the rule is meant to be used via direct calls to the provisioned endpoint.

* `description` - (Optional) The description of the Data Collection Rule.

* `identity` - (Optional) An `identity` block as defined below.

* `kind` - (Optional) The kind of the Data Collection Rule. Possible values are `Linux`, `Windows`, `AgentDirectToStore` and `WorkspaceTransforms`. A rule of kind `Linux` does not allow for `windows_event_log` data sources. And a rule of kind `Windows` does not allow for `syslog` data sources. If kind is not specified, all kinds of data sources are allowed.

~> **Note:** Once `kind` has been set, changing it forces a new Data Collection Rule to be created.

* `stream_declaration` - (Optional) A `stream_declaration` block as defined below.

* `tags` - (Optional) A mapping of tags which should be assigned to the Data Collection Rule.

---

A `azure_monitor_metrics` block supports the following:

* `name` - (Required) The name which should be used for this destination. This name should be unique across all destinations regardless of type within the Data Collection Rule.

---

A `column` block supports the following:

* `name` - (Required) The name of the column.

* `type` - (Required) The type of the column data. Possible values are `string`, `int`, `long`, `real`, `boolean`, `datetime`,and `dynamic`.

---

A `data_import` block supports the following:

* `event_hub_data_source` - (Required) An `event_hub_data_source` block as defined below.

---

A `data_flow` block supports the following:

* `destinations` - (Required) Specifies a list of destination names. A `azure_monitor_metrics` data source only allows for stream of kind `Microsoft-InsightsMetrics`.

* `streams` - (Required) Specifies a list of streams. Possible values include but not limited to `Microsoft-Event`, `Microsoft-InsightsMetrics`, `Microsoft-Perf`, `Microsoft-Syslog`, `Microsoft-WindowsEvent`, and `Microsoft-PrometheusMetrics`.

* `built_in_transform` - (Optional) The built-in transform to transform stream data.

* `output_stream` - (Optional) The output stream of the transform. Only required if the data flow changes data to a different stream.

* `transform_kql` - (Optional) The KQL query to transform stream data.

---

A `data_sources` block supports the following:

* `data_import` - (Optional) A `data_import` block as defined above.

* `extension` - (Optional) One or more `extension` blocks as defined below.

* `iis_log` - (Optional) One or more `iis_log` blocks as defined below.

* `log_file` - (Optional) One or more `log_file` blocks as defined below.

* `performance_counter` - (Optional) One or more `performance_counter` blocks as defined below.

* `platform_telemetry` - (Optional) One or more `platform_telemetry` blocks as defined below.

* `prometheus_forwarder` - (Optional) One or more `prometheus_forwarder` blocks as defined below.

* `syslog` - (Optional) One or more `syslog` blocks as defined below.

* `windows_event_log` - (Optional) One or more `windows_event_log` blocks as defined below.

* `windows_firewall_log` - (Optional) One or more `windows_firewall_log` blocks as defined below.

---

A `destinations` block supports the following:

* `azure_monitor_metrics` - (Optional) A `azure_monitor_metrics` block as defined above.

* `event_hub` - (Optional) One or more `event_hub` blocks as defined below.

* `event_hub_direct` - (Optional) One or more `event_hub` blocks as defined below.

* `log_analytics` - (Optional) One or more `log_analytics` blocks as defined below.

* `monitor_account` - (Optional) One or more `monitor_account` blocks as defined below.

* `storage_blob` - (Optional) One or more `storage_blob` blocks as defined below.

* `storage_blob_direct` - (Optional) One or more `storage_blob_direct` blocks as defined below.

* `storage_table_direct` - (Optional) One or more `storage_table_direct` blocks as defined below.

-> **Note:** `event_hub_direct`, `storage_blob_direct`, and `storage_table_direct` are only available for rules of kind `AgentDirectToStore`.

-> **Note:** At least one of `azure_monitor_metrics`, `event_hub`, `event_hub_direct`, `log_analytics`, `monitor_account`, `storage_blob`, `storage_blob_direct`,and `storage_table_direct` blocks must be specified.

---

An `event_hub_data_source` block supports the following:

* `name` - (Required) The name which should be used for this data source. This name should be unique across all data sources regardless of type within the Data Collection Rule.

* `stream` - (Required) The stream to collect from Event Hub. Possible value should be a custom stream name.

* `consumer_group` - (Optional) The Event Hub consumer group name.

---

An `event_hub` block supports the following:

* `event_hub_id` - (Required) The resource ID of the Event Hub.

* `name` - (Required) The name which should be used for this destination. This name should be unique across all destinations regardless of type within the Data Collection Rule.

---

An `event_hub_direct` block supports the following:

* `event_hub_id` - (Required) The resource ID of the Event Hub.

* `name` - (Required) The name which should be used for this destination. This name should be unique across all destinations regardless of type within the Data Collection Rule.

---

An `extension` block supports the following:

* `extension_name` - (Required) The name of the VM extension.

* `name` - (Required) The name which should be used for this data source. This name should be unique across all data sources regardless of type within the Data Collection Rule.

* `streams` - (Required) Specifies a list of streams that this data source will be sent to. A stream indicates what schema will be used for this data and usually what table in Log Analytics the data will be sent to. Possible values include but not limited to `Microsoft-Event`, `Microsoft-InsightsMetrics`, `Microsoft-Perf`, `Microsoft-Syslog`, `Microsoft-WindowsEvent`.

* `extension_json` - (Optional) A JSON String which specifies the extension setting.

* `input_data_sources` - (Optional) Specifies a list of data sources this extension needs data from. An item should be a name of a supported data source which produces only one stream. Supported data sources type: `performance_counter`, `windows_event_log`,and `syslog`.

---

An `iis_log` block supports the following:

* `name` - (Required) The name which should be used for this data source. This name should be unique across all data sources regardless of type within the Data Collection Rule.

* `streams` - (Required) Specifies a list of streams that this data source will be sent to. A stream indicates what schema will be used for this data and usually what table in Log Analytics the data will be sent to. Possible value is `Microsoft-W3CIISLog`.

* `log_directories` - (Optional) Specifies a list of absolute paths where the log files are located.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Data Collection Rule. Possible values are `SystemAssigned` and `UserAssigned`.

* `identity_ids` - (Optional) A list of User Assigned Managed Identity IDs to be assigned to this Data Collection Rule. Currently, up to 1 identity is supported.

~> **Note:** This is required when `type` is set to `UserAssigned`.

---

A `label_include_filter` block supports the following:

* `label` - (Required) The label of the filter. This label should be unique across all `label_include_fileter` block. Possible value is `microsoft_metrics_include_label`.

* `value` - (Required) The value of the filter.

---

A `log_analytics` block supports the following:

* `name` - (Required) The name which should be used for this destination. This name should be unique across all destinations regardless of type within the Data Collection Rule.

* `workspace_resource_id` - (Required) The ID of a Log Analytic Workspace resource.

---

An `log_file` block supports the following:

* `name` - (Required) The name which should be used for this data source. This name should be unique across all data sources regardless of type within the Data Collection Rule.

* `streams` - (Required) Specifies a list of streams that this data source will be sent to. A stream indicates what schema will be used for this data and usually what table in Log Analytics the data will be sent to. Possible value should be custom stream names.

* `file_patterns` - (Required) Specifies a list of file patterns where the log files are located. For example, `C:\\JavaLogs\\*.log`.

* `format` - (Required) The data format of the log files. Possible values are `text` and `json`.

* `settings` - (Optional) A `settings` block as defined below.

---

A `monitor_account` block supports the following:

* `monitor_account_id` - (Required) The resource ID of the Monitor Account.

* `name` - (Required) The name which should be used for this destination. This name should be unique across all destinations regardless of type within the Data Collection Rule.

---

A `performance_counter` block supports the following:

* `counter_specifiers` - (Required) Specifies a list of specifier names of the performance counters you want to collect. To get a list of performance counters on Windows, run the command `typeperf`. Please see [this document](https://learn.microsoft.com/en-us/azure/azure-monitor/agents/data-sources-performance-counters#configure-performance-counters) for more information.

* `name` - (Required) The name which should be used for this data source. This name should be unique across all data sources regardless of type within the Data Collection Rule.

* `sampling_frequency_in_seconds` - (Required) The number of seconds between consecutive counter measurements (samples). The value should be integer between `1` and `1800` inclusive. `sampling_frequency_in_seconds` must be equal to `60` seconds for counters collected with `Microsoft-InsightsMetrics` stream.

* `streams` - (Required) Specifies a list of streams that this data source will be sent to. A stream indicates what schema will be used for this data and usually what table in Log Analytics the data will be sent to. Possible values include but not limited to `Microsoft-InsightsMetrics`,and `Microsoft-Perf`.

---

A `platform_telemetry` block supports the following:

* `name` - (Required) The name which should be used for this data source. This name should be unique across all data sources regardless of type within the Data Collection Rule.

* `streams` - (Required) Specifies a list of streams that this data source will be sent to. A stream indicates what schema will be used for this data and usually what table in Log Analytics the data will be sent to. Possible values include but not limited to `Microsoft.Cache/redis:Metrics-Group-All`.

---

A `prometheus_forwarder` block supports the following:

* `name` - (Required) The name which should be used for this data source. This name should be unique across all data sources regardless of type within the Data Collection Rule.

* `streams` - (Required) Specifies a list of streams that this data source will be sent to. A stream indicates what schema will be used for this data and usually what table in Log Analytics the data will be sent to. Possible value is `Microsoft-PrometheusMetrics`.

* `label_include_filter` - (Optional) One or more `label_include_filter` blocks as defined above.

---

A `settings` block within the `log_file` block supports the following:

* `text` - (Required) A `text` block as defined below.

---

A `storage_blob` block supports the following:

* `container_name` - (Required) The Storage Container name.

* `name` - (Required) The name which should be used for this destination. This name should be unique across all destinations regardless of type within the Data Collection Rule.

* `storage_account_id` - (Required) The resource ID of the Storage Account.

---

A `storage_blob_direct` block supports the following:

* `container_name` - (Required) The Storage Container name.

* `name` - (Required) The name which should be used for this destination. This name should be unique across all destinations regardless of type within the Data Collection Rule.

* `storage_account_id` - (Required) The resource ID of the Storage Account.

---

---

A `storage_table_direct` block supports the following:

* `table_name` - (Required) The Storage Table name.

* `name` - (Required) The name which should be used for this destination. This name should be unique across all destinations regardless of type within the Data Collection Rule.

* `storage_account_id` - (Required) The resource ID of the Storage Account.

---

A `stream_declaration` block supports the following:

* `stream_name` - (Required) The name of the custom stream. This name should be unique across all `stream_declaration` blocks and must begin with a prefix of `Custom-`.

* `column` - (Required) One or more `column` blocks as defined above.

---

A `syslog` block supports the following:

* `facility_names` - (Required) Specifies a list of facility names. Use a wildcard `*` to collect logs for all facility names. Possible values are `alert`, `*`, `audit`, `auth`, `authpriv`, `clock`, `cron`, `daemon`, `ftp`, `kern`, `local5`, `local4`, `local1`, `local7`, `local6`, `local3`, `local2`, `local0`, `lpr`, `mail`, `mark`, `news`, `nopri`, `ntp`, `syslog`, `user` and `uucp`.

* `streams` - (Required) Specifies a list of streams that this data source will be sent to. A stream indicates what schema will be used for this data and usually what table in Log Analytics the data will be sent to. Possible values include but not limited to `Microsoft-Syslog`,and `Microsoft-CiscoAsa`, and `Microsoft-CommonSecurityLog`.

* `log_levels` - (Required) Specifies a list of log levels. Use a wildcard `*` to collect logs for all log levels. Possible values are `Debug`, `Info`, `Notice`, `Warning`, `Error`, `Critical`, `Alert`, `Emergency`,and `*`.

* `name` - (Required) The name which should be used for this data source. This name should be unique across all data sources regardless of type within the Data Collection Rule.

---

A `text` block within the `log_file.settings` block supports the following:

* `record_start_timestamp_format` - (Required) The timestamp format of the text log files. Possible values are `ISO 8601`, `YYYY-MM-DD HH:MM:SS`, `M/D/YYYY HH:MM:SS AM/PM`, `Mon DD, YYYY HH:MM:SS`, `yyMMdd HH:mm:ss`, `ddMMyy HH:mm:ss`, `MMM d hh:mm:ss`, `dd/MMM/yyyy:HH:mm:ss zzz`,and `yyyy-MM-ddTHH:mm:ssK`.

---

A `windows_event_log` block supports the following:

* `name` - (Required) The name which should be used for this data source. This name should be unique across all data sources regardless of type within the Data Collection Rule.

* `streams` - (Required) Specifies a list of streams that this data source will be sent to. A stream indicates what schema will be used for this data and usually what table in Log Analytics the data will be sent to. Possible values include but not limited to `Microsoft-Event`,and `Microsoft-WindowsEvent` and `Microsoft-SecurityEvent`.

* `x_path_queries` - (Required) Specifies a list of Windows Event Log queries in XPath expression. Please see [this document](https://learn.microsoft.com/en-us/azure/azure-monitor/agents/data-collection-rule-azure-monitor-agent?tabs=cli#filter-events-using-xpath-queries) for more information.

---

A `windows_firewall_log` block supports the following:

* `name` - (Required) The name which should be used for this data source. This name should be unique across all data sources regardless of type within the Data Collection Rule.

* `streams` - (Required) Specifies a list of streams that this data source will be sent to. A stream indicates what schema will be used for this data and usually what table in Log Analytics the data will be sent to.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Data Collection Rule.

* `immutable_id` - The immutable ID of the Data Collection Rule.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Data Collection Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the Data Collection Rule.
* `update` - (Defaults to 30 minutes) Used when updating the Data Collection Rule.
* `delete` - (Defaults to 30 minutes) Used when deleting the Data Collection Rule.

## Import

Data Collection Rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_monitor_data_collection_rule.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Insights/dataCollectionRules/rule1
```
