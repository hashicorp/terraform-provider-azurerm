---
subcategory: "Sentinel"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_sentinel_alert_rule_template"
description: |-
  Gets information about an existing Sentinel Alert Rule Template.
---

# Data Source: azurerm_sentinel_alert_rule_template

Use this data source to access information about an existing Sentinel Alert Rule Template.

## Example Usage

```hcl
data "azurerm_sentinel_alert_rule_template" "example" {
  log_analytics_workspace_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.OperationalInsights/workspaces/workspace1"
  display_name               = "Create incidents based on Azure Security Center for IoT alerts"
}

output "id" {
  value = data.azurerm_sentinel_alert_rule_template.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `log_analytics_workspace_id` - (Required) The ID of the Log Analytics Workspace.

---

* `name` - (Optional) The name of this Sentinel Alert Rule Template. Either `display_name` or `name` have to be specified.

* `display_name` - (Optional) The display name of this Sentinel Alert Rule Template. Either `display_name` or `name` have to be specified.

~> **Note:** As `display_name` is not unique, errors may occur when there are multiple Sentinel Alert Rule Template with same `display_name`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Sentinel.

* `nrt_template` - A `nrt_template` block as defined below. This only applies to Sentinel NRT Alert Rule Template.

* `security_incident_template` - A `security_incident_template` block as defined below. This only applies to Sentinel MS Security Incident Alert Rule Template.

* `scheduled_template` - A `scheduled_template` block as defined below. This only applies to Sentinel Scheduled Alert Rule Template.

---

A `nrt_template` block exports the following:

* `description` - The description of this Sentinel NRT Alert Rule Template.

* `query` - The query of this Sentinel NRT Alert Rule Template.

* `severity` - The alert severity of this Sentinel NRT Alert Rule Template.

* `tactics` - A list of categories of attacks by which to classify the rule.

---

A `security_incident_template` block exports the following:

* `description` - The description of this Sentinel MS Security Incident Alert Rule Template.

* `product_filter` - The Microsoft Security Service from where the alert will be generated.

---

A `scheduled_template` block exports the following:

* `description` - The description of this Sentinel Scheduled Alert Rule Template.

* `query` - The query of this Sentinel Scheduled Alert Rule Template.

* `query_frequency` - The ISO 8601 timespan duration between two consecutive queries.

* `query_period` - The ISO 8601 timespan duration, which determine the time period of the data covered by the query.

* `severity` - The alert severity of this Sentinel Scheduled Alert Rule Template.

* `tactics` - A list of categories of attacks by which to classify the rule.

* `trigger_operator` - The alert trigger operator, combined with `trigger_threshold`, setting alert threshold of this Sentinel Scheduled Alert Rule Template.

* `trigger_threshold` - The baseline number of query results generated, combined with `trigger_operator`, setting alert threshold of this Sentinel Scheduled Alert Rule Template.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Sentinel.
