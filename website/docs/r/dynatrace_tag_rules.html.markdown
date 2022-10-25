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

resource "azurerm_dynatrace_monitors" "example" {
  name                            = "exmpledynatracemonitor"
  resource_group_name             = azurerm_resource_group.example.name
  location                        = azurerm_resource_group.test.location
  identity_type                   = "SystemAssigned"
  monitoring_status               = "Enabled"
  marketplace_subscription_status = "Active"

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
  name       = "examplestreamanalyticscluster"
  monitor_id = azurerm_dynatrace_monitors.test.id

  log_rule {
    filtering_tag {
      name   = "Environment"
      value  = "Prod"
      action = "Include"
    }
    send_aad_logs          = true
    send_activity_logs     = true
    send_subscription_logs = true
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

* `name` - (Required) Name of the Dynatrace tag rules. Changing this forces a new resource to be created.

* `monitor_id` - (Required) Name of the Dynatrace monitor. Changing this forces a new resource to be created.

* `log_rule` - (Optional) Set of rules for sending logs for the Monitor resource. Changing this forces a new resource to be created.

* `metric_rule` - (Optional) Set of rules for sending metrics for the Monitor resource. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Dynatrace tag rules.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the Dynatrace tag rules.
* `read` - (Defaults to 5 minutes) Used when retrieving the Dynatrace tag rules.
* `update` - (Defaults to 1 hour) Used when updating the Dynatrace tag rules.
* `delete` - (Defaults to 1 hour) Used when deleting the Dynatrace tag rules.

## Import

Dynatrace tag rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_dynatrace_tag_rules.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Dynatrace.Observability/monitors/monitor1/tagRules/tagRules1
```
