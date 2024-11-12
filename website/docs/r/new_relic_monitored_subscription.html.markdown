---
subcategory: "New Relic"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_new_relic_monitored_subscription"
description: |-
  Manages an Azure Native New Relic Monitored Subscription.
---

# azurerm_new_relic_monitored_subscription

Manages an Azure Native New Relic Monitored Subscription.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "East US"
}

resource "azurerm_new_relic_monitor" "example" {
  name                = "example-nrm"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  plan {
    effective_date = "2023-06-06T00:00:00Z"
  }

  user {
    email        = "user@example.com"
    first_name   = "Example"
    last_name    = "User"
    phone_number = "+12313803556"
  }
}

data "azurerm_subscription" "another" {
  subscription_id = "00000000-0000-0000-0000-000000000000"
}

resource "azurerm_new_relic_monitored_subscription" "example" {
  monitor_id = azurerm_new_relic_monitor.example.id
  monitored_subcription {
    subscription_id                    = data.azurerm_subscription.another.subscription_id
    azure_active_directory_log_enabled = true
    activity_log_enabled               = true
    metric_enabled                     = true
    subscription_log_enabled           = true

    log_tag_filter {
      name   = "log1"
      action = "Include"
      value  = "log1"
    }

    log_tag_filter {
      name   = "log2"
      action = "Exclude"
      value  = ""
    }

    metric_tag_filter {
      name   = "metric1"
      action = "Include"
      value  = "metric1"
    }

    metric_tag_filter {
      name   = "metric2"
      action = "Exclude"
      value  = ""
    }
  }
}
```

## Arguments Reference

The following arguments are supported:

* `monitor_id` - (Required) Specifies the ID of the New Relic Monitor. Changing this forces a new Azure Native New Relic Monitored Subscription to be created.

---

* `monitored_subscription` - (Optional) One or more `monitored_subscription` blocks as defined below.

---

A `log_tag_filter` block supports the following:

* `name` - (Required) Specifies the name (also known as the key) of the tag.

* `action` - (Required) Specifies the valid actions for the tag. Possible values are `Exclude` and `Include`. Exclusion takes priority over inclusion.

* `value` - (Required) Specifies the value of the tag.

---

A `metric_tag_filter` block supports the following:

* `name` - (Required) Specifies the name (also known as the key) of the tag.

* `action` - (Required) Specifies the valid actions for the tag. Possible values are `Exclude` and `Include`. Exclusion takes priority over inclusion.

* `value` - (Required) Specifies the value of the tag.

---

A `monitored_subscription` block supports the following:

* `subscription_id` - (Required) Specifies the UUID of the subscription to be monitored.

* `activity_log_enabled` - (Optional) Whether activity logs from Azure resources should be sent to the Monitor resource. Defaults to `false`.

* `azure_active_directory_log_enabled` - (Optional) Whether Azure Active Directory logs should be sent to the Monitor resource. Defaults to `false`.

* `log_tag_filter` - (Optional) One or more `log_tag_filter` blocks as defined below.

* `metric_enabled` - (Optional) Whether metrics should be sent to the Monitor resource. Defaults to `false`.

* `metric_tag_filter` - (Optional) One or more `metric_tag_filter` blocks as defined below.

* `subscription_log_enabled` - (Optional) Whether subscription logs should be sent to the Monitor resource. Defaults to `false`.


## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Azure Native New Relic Monitored Subscription.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Azure Native New Relic Monitored Subscription.
* `read` - (Defaults to 5 minutes) Used when retrieving the Azure Native New Relic Monitored Subscription.
* `update` - (Defaults to 30 minutes) Used when updating the Azure Native New Relic Monitored Subscription.
* `delete` - (Defaults to 30 minutes) Used when deleting the Azure Native New Relic Monitored Subscription.

## Import

Azure Native New Relic Monitored Subscriptions can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_new_relic_monitored_subscription.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/NewRelic.Observability/monitors/monitor1/monitoredSubscriptions/default
```
