---
subcategory: "Monitor"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_monitor_scheduled_query_rules_log"
description: |-
  Get information about the specified LogToMetricAction Scheduled Query Rules resource.
---

# Data Source: azurerm_monitor_scheduled_query_rules_log

Use this data source to access the properties of a LogToMetricAction scheduled query rule.

## Example Usage

```hcl
data "azurerm_monitor_scheduled_query_rules_log" "example" {
  resource_group_name = "terraform-example-rg"
  name                = "tfex-queryrule"
}

output "query_rule_id" {
  value = data.azurerm_monitor_scheduled_query_rules_log.example.id
}
```

## Argument Reference

* `name` - (Required) Specifies the name of the scheduled query rule.
* `resource_group_name` - (Required) Specifies the name of the resource group where the scheduled query rule is located.

## Attributes Reference

* `id` - The ID of the scheduled query rule.
* `location` - Specifies the Azure Region where the resource should exist.
* `authorized_resource_ids` - A list of IDs of Resources referred into query.
* `criteria` - A `criteria` block as defined below.
* `data_source_id` - The resource URI over which log search query is to be run.
* `description` - The description of the scheduled query rule.
* `enabled` - Whether this scheduled query rule is enabled.

---

`criteria` supports the following:

* `dimension` - A `dimension` block as defined below.
* `metric_name` - Name of the metric.

---

`dimension` supports the following:

* `name` - Name of the dimension.
* `operator` - Operator for dimension values.
* `values` - List of dimension values.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the App Service Environment.
