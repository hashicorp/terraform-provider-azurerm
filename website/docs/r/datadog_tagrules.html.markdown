---
subcategory: "Datadog"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_datadog_monitor_tagrules"
description: |-
  Manages TagRules on the datadog Monitor.
---

# azurerm_datadog_monitor_tagrules

Manages TagRules on the datadog Monitor.

## Example Usage

### Adding TagRules on monitor
```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-datadog"
  location = "West US 2"
}

resource "azurerm_datadog_monitor_tagrules" "example" {
  name                = "example-monitor"
  resource_group_name = azurerm_resource_group.example.name
  log_rules{
    send_subscription_logs = true
  }
  metric_rules{
    filtering_tag {
        name = "Test"
        value = "Logs"
        action = "Include"
    }
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this datadog Monitor.

* `resource_group_name` - (Required) The name of the Resource Group where the datadog Monitor should exist.

---

* `rule_set_name` - (Optional) The name of the Tag Rules configuration.

* `log_rules` - (Optional) A `log_rules` block as defined below.

* `metric_rules` - (Optional) A `metric_rules` block as defined below.

---

An `log_rules` block supports the following:

* `send_aad_logs` - (Optional) Boolean flag specifying if AAD logs should be sent for the Monitor resource.

* `send_subscription_logs` - (Optional) Boolean flag specifying if Azure subscription logs should be sent for the Monitor resource.

* `send_resource_logs` - (Optional) Boolean flag specifying if Azure resource logs should be sent for the Monitor resource.

* `filtering_tag` - (Optional) A `filtering_tag` block as defined below.

> **NOTE:** List of filtering tags to be used for capturing logs. This only takes effect if SendResourceLogs flag is enabled. If empty, all resources will be captured. If only Exclude action is specified, the rules will apply to the list of all available resources. If Include actions are specified, the rules will only include resources with the associated tags.

---

A `metric_rules` block supports the following:

* `filtering_tag` - (Optional) A `filtering_tag` block as defined below.

> **NOTE:** List of filtering tags to be used for capturing metrics. If empty, all resources will be captured. If only Exclude action is specified, the rules will apply to the list of all available resources. If Include actions are specified, the rules will only include resources with the associated tags.

---

A `filtering_tag` block supports the following:

* `name` - (Required) Name of the Tag.

* `value` - (Required) Value of the Tag.

* `action` - (Required) Allowed values Enable or Disbale.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Tagrules on datadog monitor.

* `type` - The type of the monitor resource.

* `provisioning_state` - The state of Datadog monitor.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Tagrules on the datadog Monitor.
* `read` - (Defaults to 5 minutes) Used when retrieving the Tagrules on the datadog Monitor.
* `update` - (Defaults to 30 minutes) Used when updating the Tagrules on the datadog Monitor.
* `delete` - (Defaults to 30 minutes) Used when deleting the Tagrules on the datadog Monitor.

