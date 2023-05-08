---
subcategory: "Alerts Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_monitor_alert_prometheus_rule_group"
description: |-
  Manages an Alerts Management Prometheus Rule Groups.
---

# azurerm_monitor_alert_prometheus_rule_group

Manages an Alerts Management Prometheus Rule Groups.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_monitor_alert_prometheus_rule_group" "example" {
  name                = "example-amprg"
  resource_group_name = azurerm_resource_group.example.name
  location            = "West Europe"
  cluster_name        = ""
  description         = ""
  enabled             = false
  interval            = ""
  scopes              = []
  rules {
    alert      = ""
    enabled    = false
    expression = ""
    for        = ""
    record     = ""
    severity   = 0
    actions {
      action_group_id = ""
      action_properties = {
        key = ""
      }
    }
    resolve_configuration {
      auto_resolved   = false
      time_to_resolve = ""
    }
    annotations = {
      key = ""
    }
    labels = {
      key = ""
    }
  }
  tags = {
    key = "value"
  }

}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Alerts Management Prometheus Rule Groups. Changing this forces a new Alerts Management Prometheus Rule Groups to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the Alerts Management Prometheus Rule Groups should exist. Changing this forces a new Alerts Management Prometheus Rule Groups to be created.

* `location` - (Required) Specifies the Azure Region where the Alerts Management Prometheus Rule Groups should exist. Changing this forces a new Alerts Management Prometheus Rule Groups to be created.

* `rules` - (Required) A `rules` block as defined below.

* `scopes` - (Required) Target Azure Monitor workspaces resource ids. This api-version is currently limited to creating with one scope. This may change in future.

* `cluster_name` - (Optional) Apply rule to data from a specific cluster.

* `description` - (Optional) Rule group description.

* `enabled` - (Optional) Enable/disable rule group.

* `interval` - (Optional) Specifies the interval in which to run the Prometheus rule group represented in ISO 8601 duration format. Should be between 1 and 15 minutes.

* `tags` - (Optional) A mapping of tags which should be assigned to the Alerts Management Prometheus Rule Groups.

---

A `rules` block supports the following:

* `actions` - (Optional) An `actions` block as defined below.

* `alert` - (Optional) Alert rule name.

* `annotations` - (Optional) Specifies the annotations clause specifies a set of informational labels that can be used to store longer additional information such as alert descriptions or runbook links. The annotation values can be templated.

* `enabled` - (Optional) Enable/disable rule.

* `expression` - (Required) Specifies the PromQL expression to evaluate. https://prometheus.io/docs/prometheus/latest/querying/basics/. Evaluated periodically as given by 'interval', and the result recorded as a new set of time series with the metric name as given by 'record'.

* `for` - (Optional) Specifies the amount of time alert must be active before firing.

* `labels` - (Optional) Labels to add or overwrite before storing the result.

* `record` - (Optional) Recorded metrics name.

* `resolve_configuration` - (Optional) A `resolve_configuration` block as defined below.

* `severity` - (Optional) Specifies the severity of the alerts fired by the rule. Must be between 0 and 4.

---

An `actions` block supports the following:

* `action_group_id` - (Optional) Specifies the resource id of the action group to use.

* `action_properties` - (Optional) Specifies the properties of an action group object.

---

A `resolve_configuration` block supports the following:

* `auto_resolved` - (Optional) Enable alert auto-resolution.

* `time_to_resolve` - (Optional) Alert auto-resolution timeout.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Alerts Management Prometheus Rule Groups.



## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Alerts Management Prometheus Rule Groups.
* `read` - (Defaults to 5 minutes) Used when retrieving the Alerts Management Prometheus Rule Groups.
* `update` - (Defaults to 30 minutes) Used when updating the Alerts Management Prometheus Rule Groups.
* `delete` - (Defaults to 30 minutes) Used when deleting the Alerts Management Prometheus Rule Groups.

## Import

Alerts Management Prometheus Rule Groups can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_monitor_alert_prometheus_rule_group.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.AlertsManagement/prometheusRuleGroups/ruleGroup1
```
