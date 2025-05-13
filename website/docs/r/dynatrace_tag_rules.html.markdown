---
subcategory: "Dynatrace"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dynatrace_tag_rules"
description: |-
  Manages Dynatrace tag rules.
---

# azurerm_dynatrace_tag_rules

Manages Dynatrace tag rules.

## Example Usage

```hcl

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_dynatrace_monitor" "example" {
  name                            = "exmpledynatracemonitor"
  resource_group_name             = azurerm_resource_group.example.name
  location                        = azurerm_resource_group.test.location
  monitoring_enabled              = true
  marketplace_subscription_status = "Active"

  identity {
    type = "SystemAssigned"
  }

  user {
    first_name   = "Alice"
    last_name    = "Bobab"
    email        = "alice@microsoft.com"
    phone_number = "123456"
    country      = "westus"
  }

  plan {
    usage_type     = "COMMITTED"
    billing_cycle  = "MONTHLY"
    plan           = "azureportalintegration_privatepreview@TIDhjdtn7tfnxcy"
    effective_date = "2019-08-30T15:14:33Z"
  }
}

resource "azurerm_dynatrace_tag_rules" "example" {
  name       = "default"
  monitor_id = azurerm_dynatrace_monitors.test.id

  log_rule {
    filtering_tag {
      name   = "Environment"
      value  = "Prod"
      action = "Include"
    }
    send_azure_active_directory_logs_enabled = true
    send_activity_logs_enabled               = true
    send_subscription_logs_enabled           = true
  }

  metric_rule {
    filtering_tag {
      name   = "Environment"
      value  = "Prod"
      action = "Include"
    }
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Name of the Dynatrace tag rules. Currently, the only supported value is `default`. Changing this forces a new resource to be created.

* `monitor_id` - (Required) Name of the Dynatrace monitor. Changing this forces a new resource to be created.

* `log_rule` - (Optional) Set of rules for sending logs for the Monitor resource. A `log_rule` block as defined below.

* `metric_rule` - (Optional) Set of rules for sending metrics for the Monitor resource. A `metric_rule` block as defined below.

---

The `log_rule` block supports the following:

* `send_azure_active_directory_logs_enabled` - (Optional) Send Azure Active Directory logs. The default value is `false`.

* `send_activity_logs_enabled` - (Optional) Send Activity logs. The default value is `false`.

* `send_subscription_logs_enabled` - (Optional) Send Subscription logs. The default value is `false`.

* `filtering_tag` - (Optional) Filtering tag for the log rule. A `filtering_tag` block as defined below.

---

The `metric_rule` block supports the following:

* `filtering_tag` - (Optional) Filtering tag for the metric rule. A `filtering_tag` block as defined below.

---

The `filtering_tag` block supports the following:

* `name` - (Required) Name of the filtering tag.

* `value` - (Required) Value of the filtering tag.

* `action` - (Required) Action of the filtering tag. Possible values are `Include` and `Exclude`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Dynatrace tag rules.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Dynatrace tag rules.
* `read` - (Defaults to 5 minutes) Used when retrieving the Dynatrace tag rules.
* `update` - (Defaults to 30 minutes) Used when updating the Dynatrace tag rules.
* `delete` - (Defaults to 30 minutes) Used when deleting the Dynatrace tag rules.

## Import

Dynatrace tag rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_dynatrace_tag_rules.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Dynatrace.Observability/monitors/monitor1/tagRules/tagRules1
```
