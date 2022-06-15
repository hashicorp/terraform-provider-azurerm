---
subcategory: "Logz"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_logz_tag_rule"
description: |-
  Manages a logz Tag Rule.
---

# azurerm_logz_tag_rule

Manages a logz Tag Rule.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-logz"
  location = "West Europe"
}

resource "azurerm_logz_monitor" "example" {
  name                = "example-monitor"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  plan {
    billing_cycle  = "MONTHLY"
    effective_date = "2022-06-06T00:00:00Z"
    plan_id        = "100gb14days"
    usage_type     = "COMMITTED"
  }

  user {
    email        = "user@example.com"
    first_name   = "Example"
    last_name    = "User"
    phone_number = "+12313803556"
  }
}

resource "azurerm_logz_tag_rule" "example" {
  logz_monitor_id = azurerm_logz_monitor.example.id
  tag_filter {
    name   = "name1"
    action = "Include"
    value  = "value1"
  }

  tag_filter {
    name   = "name2"
    action = "Exclude"
    value  = "value2"
  }
  send_aad_logs          = true
  send_activity_logs     = true
  send_subscription_logs = true
}
```

## Arguments Reference

The following arguments are supported:

* `logz_monitor_id` - (Required) The ID of the Logz Monitor. Changing this forces a new logz Tag Rule to be created.

---

* `tag_filter` - (Optional) One or more (up to 10) `tag_filter` blocks as defined below.

* `send_aad_logs` - (Optional) Whether AAD logs should be sent to the Monitor resource?

* `send_activity_logs` - (Optional) Whether activity logs from Azure resources should be sent to the Monitor resource?

* `send_subscription_logs` - (Optional) Whether subscription logs should be sent to the Monitor resource?

---

An `tag_filter` block exports the following:

* `name` - (Required) The name of this `tag_filter`.

* `action` - (Required) The action for a filtering tag. Possible values are `Include` and `Exclude` is allowed. Note that the `Exclude` takes priority over the `Include`.

* `value` - (Optional) The value of this `tag_filter`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the logz Tag Rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the logz Tag Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the logz Tag Rule.
* `update` - (Defaults to 30 minutes) Used when updating the logz Tag Rule.
* `delete` - (Defaults to 30 minutes) Used when deleting the logz Tag Rule.

## Import

logz Tag Rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_logz_tag_rule.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Logz/monitors/monitor1/tagRules/ruleSet1
```
