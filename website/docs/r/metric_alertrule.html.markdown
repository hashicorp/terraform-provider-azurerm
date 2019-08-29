---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_metric_alertrule"
sidebar_current: "docs-azurerm-resource-monitor-metric-alertrule"
description: |-
  Manages a metric-based alert rule.

---

# azurerm_metric_alertrule

Manages a [metric-based alert rule](https://docs.microsoft.com/en-us/azure/monitoring-and-diagnostics/monitor-quick-resource-metric-alert-portal) in Azure Monitor.

~> **NOTE:** This resource has been [deprecated](https://docs.microsoft.com/en-us/azure/azure-monitor/platform/monitoring-classic-retirement) in favour of the `azurerm_monitor_metric_alert` resource and will be removed in the next major version of the AzureRM Provider. The new resource shares the same fields as this one, and information on migrating across [can be found in this guide](../guides/migrating-between-renamed-resources.html).

## Example Usage (CPU Percentage of a virtual machine)

```hcl
resource "azurerm_metric_alertrule" "test" {
  name                = "${azurerm_virtual_machine.test.name}-cpu"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  description = "An alert rule to watch the metric Percentage CPU"

  enabled = true

  resource_id = "${azurerm_virtual_machine.test.id}"
  metric_name = "Percentage CPU"
  operator    = "GreaterThan"
  threshold   = 75
  aggregation = "Average"
  period      = "PT5M"

  email_action {
    send_to_service_owners = false

    custom_emails = [
      "some.user@example.com",
    ]
  }

  webhook_action {
    service_uri = "https://example.com/some-url"

    properties = {
      severity        = "incredible"
      acceptance_test = "true"
    }
  }
}
```

## Example Usage (Storage usage of a SQL Database)

```hcl
resource "azurerm_metric_alertrule" "test" {
  name                = "${azurerm_sql_database.test.name}-storage"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  description = "An alert rule to watch the metric Storage"

  enabled = true

  resource_id = "${azurerm_sql_database.test.id}"
  metric_name = "storage"
  operator    = "GreaterThan"
  threshold   = 1073741824
  aggregation = "Maximum"
  period      = "PT10M"

  email_action {
    send_to_service_owners = false

    custom_emails = [
      "some.user@example.com",
    ]
  }

  webhook_action {
    service_uri = "https://example.com/some-url"

    properties = {
      severity        = "incredible"
      acceptance_test = "true"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the alert rule. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the alert rule. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `description` - (Optional) A verbose description of the alert rule that will be included in the alert email.

* `enabled` - (Optional) If `true`, the alert rule is enabled. Defaults to `true`.

---

* `resource_id` - (Required) The ID of the resource monitored by the alert rule.

* `metric_name` - (Required) The metric that defines what the rule monitors.

-> For a comprehensive reference of supported `metric_name` values for types of `resource` refer to [Supported metrics with Azure Monitor](https://docs.microsoft.com/en-us/azure/monitoring-and-diagnostics/monitoring-supported-metrics) in the Azure documentation. In the referred table, the column "Metric" corresponds to supported values for `metric_name`.

* `operator` - (Required) The operator used to compare the metric data and the threshold. Possible values are `GreaterThan`, `GreaterThanOrEqual`, `LessThan`, and `LessThanOrEqual`.

* `threshold` - (Required) The threshold value that activates the alert.

* `period` - (Required) The period of time formatted in [ISO 8601 duration format](https://en.wikipedia.org/wiki/ISO_8601#Durations) that is used to monitor the alert activity based on the threshold. The period must be between 5 minutes and 1 day.

* `aggregation` - (Required) Defines how the metric data is combined over time. Possible values are `Average`, `Minimum`, `Maximum`, `Total`, and `Last`.

---

* `email_action` - (Optional) A `email_action` block as defined below.

* `webhook_action` - (Optional) A `webhook_action` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource. Changing this forces a new resource to be created.

---

`email_action` supports the following:

* `send_to_service_owners` - (Optional) If `true`, the administrators (service and co-administrators) of the subscription are notified when the alert is triggered. Defaults to `false`.

* `custom_emails` - (Optional) A list of email addresses to be notified when the alert is triggered.

---

`webhook_action` supports the following:

* `service_uri` - (Required) The service uri of the webhook to POST the notification when the alert is triggered.

* `properties` - (Optional) A dictionary of custom properties to include with the webhook POST operation payload.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the alert rule.

## Import

Metric Alert Rules can be imported using the `resource id`, e.g.

```
terraform import azurerm_metric_alertrule.alertrule1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/microsoft.insights/alertrules/alertrule1
```
