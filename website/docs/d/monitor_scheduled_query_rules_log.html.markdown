---
subcategory: "Monitor"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_monitor_scheduled_query_rules_log"
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
* `resource_group_name` - (Required) Specifies the name of the resource group where the Scheduled Query Rule is located.

## Attributes Reference

* `id` - The ID of the Scheduled Query Rule.
* `criteria` - A `criteria` block as defined below.
* `data_source_id` - The resource URI over which log search query is to be run.
* `description` - The description of the Scheduled Query Rule.
* `enabled` - Whether this Scheduled Query Rule is enabled.

---

`criteria` supports the following:

* `dimension` - A `dimension` block as defined below.
* `metric_name` - Name of the metric.

---

`dimension` supports the following:

* `name` - Name of the dimension.
* `operator` - Operator for dimension values.
* `values` - List of dimension values.
