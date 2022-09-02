---
subcategory: "Datadog"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_datadog_monitor_tagrule"
description: |-
  Manages Tag Rules on the datadog Monitor.
---

# azurerm_datadog_monitor_tagrule

Manages TagRules on the datadog Monitor.

## Example Usage

### Adding TagRules on monitor
```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-datadog"
  location = "West US 2"
}
resource "azurerm_datadog_monitor_tagrule" "example" {
  datadog_monitor_id = azurerm_datadog_monitor.example.id
  log {
    subscription_log_enabled = true
  }
  metric {
    filtering_tag {
      name   = "Test"
      value  = "Logs"
      action = "Include"
    }
  }
}
```

## Arguments Reference

The following arguments are supported:

* `datadog_monitor_id` - (Required) The Datadog Monitor Id which should be used for this.

* `resource_group_name` - (Required) The name of the Resource Group where the datadog Monitor should exist.

---

* `rule_set_name` - (Optional) The name of the Tag Rules configuration.

* `log` - (Optional) A `log_rules` block as defined below.

* `metric` - (Optional) A `metric_rules` block as defined below.

---

An `log` block supports the following:

* `aad_log_enabled` - (Optional) Boolean flag specifying if AAD logs should be sent for the Monitor resource.

* `subscription_log_enabled` - (Optional) Boolean flag specifying if Azure subscription logs should be sent for the Monitor resource.

* `resource_log_enabled` - (Optional) Boolean flag specifying if Azure resource logs should be sent for the Monitor resource.

* `filtering_tag` - (Optional) A `filtering_tag` block as defined below.

> **NOTE:** List of filtering tags to be used for capturing logs. This only takes effect if SendResourceLogs flag is enabled. If empty, all resources will be captured. If only Exclude action is specified, the rules will apply to the list of all available resources. If Include actions are specified, the rules will only include resources with the associated tags.
---

A `metric` block supports the following:

* `filtering_tag` - (Optional) A `filtering_tag` block as defined below.

> **NOTE:** List of filtering tags to be used for capturing metrics. If empty, all resources will be captured. If only Exclude action is specified, the rules will apply to the list of all available resources. If Include actions are specified, the rules will only include resources with the associated tags.
---

A `filtering_tag` block supports the following:

* `name` - (Required) Name of the Tag.

* `value` - (Required) Value of the Tag.

* `action` - (Required) Allowed values Enable or Disable.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Tagrules on the datadog Monitor.
* `read` - (Defaults to 5 minutes) Used when retrieving the Tagrules on the datadog Monitor.
* `update` - (Defaults to 30 minutes) Used when updating the Tagrules on the datadog Monitor.
* `delete` - (Defaults to 30 minutes) Used when deleting the Tagrules on the datadog Monitor.
