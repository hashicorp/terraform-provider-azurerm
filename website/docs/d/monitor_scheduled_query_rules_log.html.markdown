---
subcategory: "Monitor"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_monitor_scheduled_query_rules_log"
sidebar_current: "docs-azurerm-datasource-monitor-scheduled-query-rules-log"
description: |-
  Get information about the specified LogToMetricAction Scheduled Query Rule.
---

# Data Source: azurerm_monitor_scheduled_query_rules_log

Use this data source to access the properties of a LogToMetricAction Scheduled Query Rule.

## Example Usage

```hcl
data "azurerm_monitor_scheduled_query_rules_log" "example" {
  resource_group_name = "terraform-example-rg"
  name                = "tfex-queryrule"
}

output "query_rule_id" {
  value = "${data.azurerm_monitor_scheduled_query_rules_log.example.id}"
}
```

## Argument Reference

* `name` - (Required) Specifies the name of the Scheduled Query Rule.
* `resource_group_name` - (Required) Specifies the name of the resource group the Scheduled Query Rule is located in.

## Attributes Reference

* `id` - The ID of the Scheduled Query Rule.
* `criteria` - A `criteria` block as defined below.
* `data_source_id` - The resource uri over which log search query is to be run.
* `description` - The description of the Scheduled Query Rule.
* `enabled` - Whether this scheduled query rule is enabled.
* `throttling` - Time (in minutes) for which Alerts should be throttled or suppressed.

---

`criteria` supports the following:

* `dimension` - (Required) A `dimension` block as defined below.
* `metric_name` - (Required) Name of the metric.

---

`dimension` supports the following:

* `name` - (Required) Name of the dimension.
* `operator` - (Required) Operator for dimension values, - 'Include'.
* `values` - (Required) List of dimension values.
