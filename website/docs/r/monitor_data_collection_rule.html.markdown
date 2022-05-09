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
  location = "East US 2"
}

resource "azurerm_log_analytics_workspace" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
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
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  azure_monitor_metrics_destination {
    name = "amm1"
  }

  log_analytics_destination {
    name                  = "centralWorkspace"
    workspace_resource_id = azurerm_log_analytics_workspace.example.id
  }

  windows_event_log_data_source {
    name          = "cloudSecurityTeamEvents"
    streams       = ["Microsoft-WindowsEvent"]
    xpath_queries = ["Security!"]
  }

  windows_event_log_data_source {
    name    = "appTeam1AppEvents"
    streams = ["Microsoft-WindowsEvent"]
    xpath_queries = ["System![System[(Level = 1 or Level = 2 or Level = 3)]]",
    "Application!*[System[(Level = 1 or Level = 2 or Level = 3)]]"]
  }

  syslog_data_source {
    name           = "cronSyslog"
    streams        = ["Microsoft-Syslog"]
    log_levels     = ["Debug", "Critical", "Emergency"]
    facility_names = ["cron"]
  }

  syslog_data_source {
    name           = "syslogBase"
    streams        = ["Microsoft-Syslog"]
    log_levels     = ["Alert", "Critical", "Emergency"]
    facility_names = ["syslog"]
  }

  extension_data_source {
    name               = "extension1"
    extension_name     = "mockname"
    streams            = ["Microsoft-Event"]
    input_data_sources = []
    extension_setting  = <<BODY
{
    "key1": "value1",
    "key2": "value2"
}
BODY
  }

  performance_counter_data_source {
    name    = "cloudTeamCoreCounters"
    streams = ["Microsoft-Perf"]
    specifiers = [
      "\\\\Memory\\\\Committed Bytes",
      "\\\\LogicalDisk(_Total)\\\\Free Megabytes",
      "\\\\PhysicalDisk(_Total)\\\\Avg. Disk Queue Length"
    ]
    sampling_frequency = 15
  }

  performance_counter_data_source {
    name    = "appTeamExtraCounters"
    streams = ["Microsoft-Perf"]
    specifiers = [
      "\\\\Process(_Total)\\\\Thread Count"
    ]
    sampling_frequency = 30
  }

  data_flows {
    streams      = ["Microsoft-InsightsMetrics"]
    destinations = ["amm1"]
  }

  data_flows {
    streams      = ["Microsoft-Perf", "Microsoft-Syslog", "Microsoft-WindowsEvent"]
    destinations = ["centralWorkspace"]
  }

  description = "this is description"

  tags = {
    Environment = "Production"
  }

  depends_on = [
    azurerm_log_analytics_solution.example
  ]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Data Collection Rule. Changing this forces a new Data Collection Rule to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Data Collection Rule should exist. Changing this forces a new Data Collection Rule to be created.

* `location` - (Required) The Azure Region where the Data Collection Rule should exist. Changing this forces a new Data Collection Rule to be created.

* `data_flows` - (Required) One or more `data_flows` blocks as defined below.

* `kind` - (Required) The kind of the resource. Accepted values are `Linux` and `Windows`. Default to support both `Linux` and `Windows` if omitted it.

---

* `azure_monitor_metrics_destination` - (Optional) A `azure_monitor_metrics_destination` block as defined below.

* `description` - (Optional) Description of the Data Collection Rule.

* `extension_data_source` - (Optional) One or more `extension_data_source` blocks as defined below.

* `log_analytics_destination` - (Optional) One or more `log_analytics_destination` blocks as defined below.

* `performance_counter_data_source` - (Optional) One or more `performance_counter_data_source` blocks as defined below.

* `syslog_data_source` - (Optional) One or more `syslog_data_source` blocks as defined below.

* `windows_event_log_data_source` - (Optional) One or more `windows_event_log_data_source` blocks as defined below.

* `tags` - (Optional) A mapping of tags which should be assigned to the Data Collection Rule.
---

A `azure_monitor_metrics_destination` block supports the following:

* `name` - (Required) A friendly name for the destination. This name should be unique across all destinations (regardless of type) within the data collection rule.

---

A `data_flows` block supports the following:

* `destinations` - (Required) List of destinations for this data flow.

* `streams` - (Required) List of streams for this data flow. Accepted values are `Microsoft-Event`, `Microsoft-InsightsMetrics`, `Microsoft-Perf`, `Microsoft-Syslog` and `Microsoft-WindowsEvent`.

---

A `extension_data_source` block supports the following:

* `extension_name` - (Required) The name of the VM extension.

* `name` - (Required) A friendly name for the data source. This name should be unique across all data sources (regardless of type) within the data collection rule.

* `streams` - (Required) List of streams that this data source will be sent to. A stream indicates what schema will be used for this data and usually what table in Log Analytics the data will be sent to. Accepted values are `Microsoft-Event`, `Microsoft-InsightsMetrics`, `Microsoft-Perf`, `Microsoft-Syslog` and `Microsoft-WindowsEvent`. 

* `extension_setting` - (Optional) The extension settings. The format is specific for particular extension.

* `input_data_sources` - (Optional) The list of data sources this extension needs data from.

---

A `log_analytics_destination` block supports the following:

* `name` - (Required) A friendly name for the destination. This name should be unique across all destinations (regardless of type) within the data collection rule.

* `workspace_resource_id` - (Required) The resource ID of the Log Analytics workspace.

---

A `performance_counter_data_source` block supports the following:

* `name` - (Required) The name which should be used for this TODO.

* `streams` - (Required) List of streams that this data source will be sent to. A stream indicates what schema will be used for this data and usually what table in Log Analytics the data will be sent to. Accepted values are `Microsoft-InsightsMetrics` and `Microsoft-Perf`.

* `sampling_frequency` - (Optional) The number of seconds between consecutive counter measurements (samples).

* `specifiers` - (Optional) A list of specifier names of the performance counters you want to collect. Use a wildcard (*) to collect a counter for all instances. To get a list of performance counters on Windows, run the command `typeperf`.

---

A `syslog_data_source` block supports the following:

* `name` - (Required) A friendly name for the data source. This name should be unique across all data sources (regardless of type) within the data collection rule.

* `streams` - (Required) List of streams that this data source will be sent to. A stream indicates what schema will be used for this data and usually what table in Log Analytics the data will be sent to. The accepted value is `Microsoft-Syslog`.

* `facility_names` - (Optional) The list of facility names.

* `log_levels` - (Optional) The log levels to collect.

---

A `windows_event_log_data_source` block supports the following:

* `name` - (Required) A friendly name for the data source. This name should be unique across all data sources (regardless of type) within the data collection rule.

* `streams` - (Required) List of streams that this data source will be sent to. A stream indicates what schema will be used for this data and usually what table in Log Analytics the data will be sent to. The accepted value is `Microsoft-Event` and `Microsoft-WindowsEvent`.

* `xpath_queries` - (Optional) A list of Windows Event Log queries in XPATH format.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Data Collection Rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Data Collection Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the Data Collection Rule.
* `update` - (Defaults to 30 minutes) Used when updating the Data Collection Rule.
* `delete` - (Defaults to 30 minutes) Used when deleting the Data Collection Rule.

## Import

Data Collection Rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_monitor_data_collection_rule.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Insights/dataCollectionRules/rule1
```
