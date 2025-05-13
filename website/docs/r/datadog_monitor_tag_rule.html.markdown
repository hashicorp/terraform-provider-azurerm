---
subcategory: "Datadog"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_datadog_monitor_tag_rule"
description: |-
  Manages Tag Rules on the Datadog Monitor.
---

# azurerm_datadog_monitor_tag_rule

Manages TagRules on the datadog Monitor.

## Example Usage

### Adding TagRules on monitor
```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-datadog"
  location = "West US 2"
}

resource "azurerm_datadog_monitor" "example" {
  name                = "example-monitor"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  datadog_organization {
    api_key         = "XXXX"
    application_key = "XXXX"
  }
  user {
    name  = "Example"
    email = "abc@xyz.com"
  }
  sku_name = "Linked"
  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_datadog_monitor_tag_rule" "example" {
  datadog_monitor_id = azurerm_datadog_monitor.example.id
  log {
    subscription_log_enabled = true
  }
  metric {
    filter {
      name   = "Test"
      value  = "Logs"
      action = "Include"
    }
  }
}
```

## Arguments Reference

The following arguments are supported:

* `datadog_monitor_id` - (Required) The Datadog Monitor Id which should be used for this Datadog Monitor Tag Rule. Changing this forces a new Datadog Monitor Tag Rule to be created.

---

* `name` - (Optional) The name of the Tag Rules configuration. The allowed value is `default`. Defaults to `default`.

* `log` - (Optional) A `log` block as defined below.

* `metric` - (Optional) A `metric` block as defined below.

---

An `log` block supports the following:

* `aad_log_enabled` - (Optional) Whether AAD logs should be sent for the Monitor resource?

* `subscription_log_enabled` - (Optional) Whether Azure subscription logs should be sent for the Monitor resource?

* `resource_log_enabled` - (Optional) Whether Azure resource logs should be sent for the Monitor resource?

* `filter` - (Optional) A `filter` block as defined below.

-> **Note:** List of filtering tags to be used for capturing logs. This only takes effect if `resource_log_enabled` flag is enabled. If empty, all resources will be captured. If only Exclude action is specified, the rules will apply to the list of all available resources. If Include actions are specified, the rules will only include resources with the associated tags.

---

A `metric` block supports the following:

* `filter` - (Optional) A `filter` block as defined below.

-> **Note:** List of filtering tags to be used for capturing metrics. If empty, all resources will be captured. If only Exclude action is specified, the rules will apply to the list of all available resources. If Include actions are specified, the rules will only include resources with the associated tags.

---

A `filter` block supports the following:

* `name` - (Required) Name of the Tag.

* `value` - (Required) Value of the Tag.

* `action` - (Required) Allowed values Include or Exclude.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Tagrules on the datadog Monitor.
* `read` - (Defaults to 5 minutes) Used when retrieving the Tagrules on the datadog Monitor.
* `update` - (Defaults to 30 minutes) Used when updating the Tagrules on the datadog Monitor.
* `delete` - (Defaults to 30 minutes) Used when deleting the Tagrules on the datadog Monitor.

## Import

Tag Rules on the Datadog Monitor can be imported using the `tag rule resource id`, e.g.

```shell
terraform import azurerm_datadog_monitor_tag_rule.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Datadog/monitors/monitor1/tagRules/default
