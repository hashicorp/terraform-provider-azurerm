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

resource "azurerm_monitor_data_collection_rule" "example" {
  name                = "example-rule"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  destinations {
    log_analytics {
      workspace_resource_id = azurerm_log_analytics_workspace.example.id
      name                  = "test-destination-log"
    }

    azure_monitor_metrics {
      name = "test-destination-metrics"
    }
  }

  data_flow {
    streams      = ["Microsoft-InsightsMetrics"]
    destinations = ["test-destination-metrics"]
  }

  data_flow {
    streams      = ["Microsoft-InsightsMetrics", "Microsoft-Syslog", "Microsoft-Perf"]
    destinations = ["test-destination-log"]
  }

  data_sources {
    syslog {
      facility_names = ["*"]
      log_levels     = ["*"]
      name           = "test-datasource-syslog"
    }

    performance_counter {
      streams                       = ["Microsoft-Perf", "Microsoft-InsightsMetrics"]
      sampling_frequency_in_seconds = 10
      counter_specifiers            = ["Processor(*)\\% Processor Time"]
      name                          = "test-datasource-perfcounter"
    }

    windows_event_log {
      streams        = ["Microsoft-WindowsEvent"]
      x_path_queries = ["*[System/Level=1]"]
      name           = "test-datasource-wineventlog"
    }

    extension {
      streams            = ["Microsoft-WindowsEvent"]
      input_data_sources = ["test-datasource-wineventlog"]
      extension_name     = "test-extension-name"
      extension_json = jsonencode({
        a = 1
        b = "hello"
      })
      name = "test-datasource-extension"
    }
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

* `data_sources` - (Optional) A `data_sources` block as defined below. This property is optional and can be omitted if the rule is meant to be used via direct calls to the provisioned endpoint.

* `description` - (Optional) The description of the Data Collection Rule.

* `kind` - (Optional) The kind of the Data Collection Rule. Possible values are `Linux` and `Windows`. A rule of kind `Linux` does not allow for `windows_event_log` data sources. And a rule of kind `Windows` does not allow for `syslog` data sources. If kind is not specified, all kinds of data sources are allowed.

* `tags` - (Optional) A mapping of tags which should be assigned to the Data Collection Rule.

---

A `azure_monitor_metrics` block supports the following:

* `name` - (Required) The name which should be used for this destination. This name should be unique across all destinations regardless of type within the Data Collection Rule.

---

A `data_flow` block supports the following:

* `destinations` - (Required) Specifies a list of destination names. A `azure_monitor_metrics` data source only allows for stream of kind `Microsoft-InsightsMetrics`.

* `streams` - (Required) Specifies a list of streams. Possible values include but not limited to `Microsoft-Event`, `Microsoft-InsightsMetrics`, `Microsoft-Perf`, `Microsoft-Syslog`,and `Microsoft-WindowsEvent`.

---

A `data_sources` block supports the following:

* `extension` - (Optional) One or more `extension` blocks as defined below.

* `performance_counter` - (Optional) One or more `performance_counter` blocks as defined below.

* `syslog` - (Optional) One or more `syslog` blocks as defined below.

* `windows_event_log` - (Optional) One or more `windows_event_log` blocks as defined below.

---

A `destinations` block supports the following:

* `azure_monitor_metrics` - (Optional) A `azure_monitor_metrics` block as defined above.

* `log_analytics` - (Optional) One or more `log_analytics` blocks as defined below.

-> **NOTE** At least one of `azure_monitor_metrics` and `log_analytics` blocks must be specified.

---

A `extension` block supports the following:

* `extension_name` - (Required) The name of the VM extension.

* `name` - (Required) The name which should be used for this data source. This name should be unique across all data sources regardless of type within the Data Collection Rule.

* `streams` - (Required) Specifies a list of streams that this data source will be sent to. A stream indicates what schema will be used for this data and usually what table in Log Analytics the data will be sent to. Possible values include but not limited to `Microsoft-Event`, `Microsoft-InsightsMetrics`, `Microsoft-Perf`, `Microsoft-Syslog`, `Microsoft-WindowsEvent`.

* `extension_json` - (Optional) A JSON String which specifies the extension setting.

* `input_data_sources` - (Optional) Specifies a list of data sources this extension needs data from. An item should be a name of a supported data source which produces only one stream. Supported data sources type: `performance_counter`, `windows_event_log`,and `syslog`.

---

A `log_analytics` block supports the following:

* `name` - (Required) The name which should be used for this destination. This name should be unique across all destinations regardless of type within the Data Collection Rule.

* `workspace_resource_id` - (Required) The ID of a Log Analytic Workspace resource.

---

A `performance_counter` block supports the following:

* `counter_specifiers` - (Required) Specifies a list of specifier names of the performance counters you want to collect. Use a wildcard `*` to collect counters for all instances. To get a list of performance counters on Windows, run the command `typeperf`.

* `name` - (Required) The name which should be used for this data source. This name should be unique across all data sources regardless of type within the Data Collection Rule.

* `sampling_frequency_in_seconds` - (Required) The number of seconds between consecutive counter measurements (samples). The value should be integer between `1` and `300` inclusive.

* `streams` - (Required) Specifies a list of streams that this data source will be sent to. A stream indicates what schema will be used for this data and usually what table in Log Analytics the data will be sent to. Possible values include but not limited to `Microsoft-InsightsMetrics`,and `Microsoft-Perf`.

---

A `syslog` block supports the following:

* `facility_names` - (Required) Specifies a list of facility names. Use a wildcard `*` to collect logs for all facility names. Possible values are `auth`, `authpriv`, `cron`, `daemon`, `kern`, `lpr`, `mail`, `mark`, `news`, `syslog`, `user`, `uucp`, `local0`, `local1`, `local2`, `local3`, `local4`, `local5`, `local6`, `local7`,and `*`.

* `log_levels` - (Required) Specifies a list of log levels. Use a wildcard `*` to collect logs for all log levels. Possible values are `Debug`, `Info`, `Notice`, `Warning`, `Error`, `Critical`, `Alert`, `Emergency`,and `*`.

* `name` - (Required) The name which should be used for this data source. This name should be unique across all data sources regardless of type within the Data Collection Rule.

* `streams` - (Optional) Specifies a list of streams that this data source will be sent to. A stream indicates what schema will be used for this data and usually what table in Log Analytics the data will be sent to. Possible values include but not limited to `Microsoft-Syslog`,and `Microsoft-CiscoAsa`, and `Microsoft-CommonSecurityLog`.

-> **Note:** In 4.0 or later version of the provider, `streams` will be required. In 3.x version of provider, if `streams` is not specified in creation, it is default to `["Microsoft-Syslog"]`. if `streams` need to be modified (include change other value to the default value), it must be explicitly specified.

---

A `windows_event_log` block supports the following:

* `name` - (Required) The name which should be used for this data source. This name should be unique across all data sources regardless of type within the Data Collection Rule.

* `streams` - (Required) Specifies a list of streams that this data source will be sent to. A stream indicates what schema will be used for this data and usually what table in Log Analytics the data will be sent to. Possible values include but not limited to `Microsoft-Event`,and `Microsoft-WindowsEvent`, `Microsoft-RomeDetectionEvent`, and `Microsoft-SecurityEvent`.

* `x_path_queries` - (Required) Specifies a list of Windows Event Log queries in XPath expression.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Data Collection Rule.

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
