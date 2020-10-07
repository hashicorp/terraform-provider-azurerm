---
subcategory: "Monitor"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_monitor_scheduled_query_rules_alert"
description: |-
  Get information about the specified AlertingAction Scheduled Query Rules resource.
---

# Data Source: azurerm_monitor_scheduled_query_rules_alert

Use this data source to access the properties of an AlertingAction scheduled query rule.

## Example Usage

```hcl
data "azurerm_monitor_scheduled_query_rules_alert" "example" {
  resource_group_name = "terraform-example-rg"
  name                = "tfex-queryrule"
}

output "query_rule_id" {
  value = "${data.azurerm_monitor_scheduled_query_rules_alert.example.id}"
}
```

## Argument Reference

* `name` - (Required) Specifies the name of the scheduled query rule.
* `resource_group_name` - (Required) Specifies the name of the resource group where the scheduled query rule is located.

## Attributes Reference

* `id` - The ID of the scheduled query rule.
* `action` - An `action` block as defined below.
* `authorized_resource_ids` - The list of Resource IDs referred into query.
* `data_source_id` - The resource URI over which log search query is to be run.
* `description` - The description of the scheduled query rule.
* `enabled` - Whether this scheduled query rule is enabled.
* `frequency` - Frequency at which rule condition should be evaluated.
* `query` - Log search query.
* `time_window` - Time window for which data needs to be fetched for query.
* `severity` - Severity of the alert.
* `throttling` - Time for which alerts should be throttled or suppressed.
* `trigger` - A `trigger` block as defined below.

---

* `action` supports the following:

* `action_group` - List of action group reference resource IDs.
* `custom_webhook_payload` - Custom payload to be sent for all webhook URI in Azure action group.
* `email_subject` - Custom subject override for all email IDs in Azure action group.

---

`metricTrigger` supports the following:

* `metricColumn` - Evaluation of metric on a particular column.
* `metricTriggerType` - The metric trigger type.
* `operator` - Evaluation operation for rule.
* `threshold` - The threshold of the metric trigger.

---

`trigger` supports the following:

* `metricTrigger` - A `metricTrigger` block as defined above.
* `operator` - Evaluation operation for rule.
* `threshold` - Result or count threshold based on which rule should be triggered.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the App Service Environment.
