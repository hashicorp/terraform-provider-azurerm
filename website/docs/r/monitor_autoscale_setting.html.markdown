---
subcategory: "Monitor"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_monitor_autoscale_setting"
description: |-
    Manages an AutoScale Setting which can be applied to Virtual Machine Scale Sets, App Services and other scalable resources.
---

# azurerm_monitor_autoscale_setting

Manages a AutoScale Setting which can be applied to Virtual Machine Scale Sets, App Services and other scalable resources.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "autoscalingTest"
  location = "West Europe"
}

resource "azurerm_virtual_machine_scale_set" "example" {
  # ...
}

resource "azurerm_monitor_autoscale_setting" "example" {
  name                = "myAutoscaleSetting"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  target_resource_id  = azurerm_virtual_machine_scale_set.example.id

  profile {
    name = "defaultProfile"

    capacity {
      default = 1
      minimum = 1
      maximum = 10
    }

    rule {
      metric_trigger {
        metric_name        = "Percentage CPU"
        metric_resource_id = azurerm_virtual_machine_scale_set.example.id
        time_grain         = "PT1M"
        statistic          = "Average"
        time_window        = "PT5M"
        time_aggregation   = "Average"
        operator           = "GreaterThan"
        threshold          = 75
        metric_namespace   = "microsoft.compute/virtualmachinescalesets"
        dimensions {
          name     = "AppName"
          operator = "Equals"
          values   = ["App1"]
        }
      }

      scale_action {
        direction = "Increase"
        type      = "ChangeCount"
        value     = "1"
        cooldown  = "PT1M"
      }
    }

    rule {
      metric_trigger {
        metric_name        = "Percentage CPU"
        metric_resource_id = azurerm_virtual_machine_scale_set.example.id
        time_grain         = "PT1M"
        statistic          = "Average"
        time_window        = "PT5M"
        time_aggregation   = "Average"
        operator           = "LessThan"
        threshold          = 25
      }

      scale_action {
        direction = "Decrease"
        type      = "ChangeCount"
        value     = "1"
        cooldown  = "PT1M"
      }
    }
  }

  notification {
    email {
      send_to_subscription_administrator    = true
      send_to_subscription_co_administrator = true
      custom_emails                         = ["admin@contoso.com"]
    }
  }
}
```

## Example Usage (repeating on weekends)

```hcl
resource "azurerm_resource_group" "example" {
  name     = "autoscalingTest"
  location = "West Europe"
}

resource "azurerm_virtual_machine_scale_set" "example" {
  # ...
}

resource "azurerm_monitor_autoscale_setting" "example" {
  name                = "myAutoscaleSetting"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  target_resource_id  = azurerm_virtual_machine_scale_set.example.id

  profile {
    name = "Weekends"

    capacity {
      default = 1
      minimum = 1
      maximum = 10
    }

    rule {
      metric_trigger {
        metric_name        = "Percentage CPU"
        metric_resource_id = azurerm_virtual_machine_scale_set.example.id
        time_grain         = "PT1M"
        statistic          = "Average"
        time_window        = "PT5M"
        time_aggregation   = "Average"
        operator           = "GreaterThan"
        threshold          = 90
      }

      scale_action {
        direction = "Increase"
        type      = "ChangeCount"
        value     = "2"
        cooldown  = "PT1M"
      }
    }

    rule {
      metric_trigger {
        metric_name        = "Percentage CPU"
        metric_resource_id = azurerm_virtual_machine_scale_set.example.id
        time_grain         = "PT1M"
        statistic          = "Average"
        time_window        = "PT5M"
        time_aggregation   = "Average"
        operator           = "LessThan"
        threshold          = 10
      }

      scale_action {
        direction = "Decrease"
        type      = "ChangeCount"
        value     = "2"
        cooldown  = "PT1M"
      }
    }

    recurrence {
      frequency = "Week"
      timezone  = "Pacific Standard Time"
      days      = ["Saturday", "Sunday"]
      hours     = [12]
      minutes   = [0]
    }
  }

  notification {
    email {
      send_to_subscription_administrator    = true
      send_to_subscription_co_administrator = true
      custom_emails                         = ["admin@contoso.com"]
    }
  }
}
```

## Example Usage (for fixed dates)

```hcl
resource "azurerm_resource_group" "example" {
  name     = "autoscalingTest"
  location = "West Europe"
}

resource "azurerm_virtual_machine_scale_set" "example" {
  # ...
}

resource "azurerm_monitor_autoscale_setting" "example" {
  name                = "myAutoscaleSetting"
  enabled             = true
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  target_resource_id  = azurerm_virtual_machine_scale_set.example.id

  profile {
    name = "forJuly"

    capacity {
      default = 1
      minimum = 1
      maximum = 10
    }

    rule {
      metric_trigger {
        metric_name        = "Percentage CPU"
        metric_resource_id = azurerm_virtual_machine_scale_set.example.id
        time_grain         = "PT1M"
        statistic          = "Average"
        time_window        = "PT5M"
        time_aggregation   = "Average"
        operator           = "GreaterThan"
        threshold          = 90
      }

      scale_action {
        direction = "Increase"
        type      = "ChangeCount"
        value     = "2"
        cooldown  = "PT1M"
      }
    }

    rule {
      metric_trigger {
        metric_name        = "Percentage CPU"
        metric_resource_id = azurerm_virtual_machine_scale_set.example.id
        time_grain         = "PT1M"
        statistic          = "Average"
        time_window        = "PT5M"
        time_aggregation   = "Average"
        operator           = "LessThan"
        threshold          = 10
      }

      scale_action {
        direction = "Decrease"
        type      = "ChangeCount"
        value     = "2"
        cooldown  = "PT1M"
      }
    }

    fixed_date {
      timezone = "Pacific Standard Time"
      start    = "2020-07-01T00:00:00Z"
      end      = "2020-07-31T23:59:59Z"
    }
  }

  notification {
    email {
      send_to_subscription_administrator    = true
      send_to_subscription_co_administrator = true
      custom_emails                         = ["admin@contoso.com"]
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the AutoScale Setting. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in the AutoScale Setting should be created. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the AutoScale Setting should exist. Changing this forces a new resource to be created.

* `profile` - (Required) Specifies one or more (up to 20) `profile` blocks as defined below.

* `target_resource_id` - (Required) Specifies the resource ID of the resource that the autoscale setting should be added to.

* `enabled` - (Optional) Specifies whether automatic scaling is enabled for the target resource. Defaults to `true`.

* `notification` - (Optional) Specifies a `notification` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `profile` block supports the following:

* `name` - (Required) Specifies the name of the profile.

* `capacity` - (Required) A `capacity` block as defined below.

* `rule` - (Optional) One or more (up to 10) `rule` blocks as defined below.

* `fixed_date` - (Optional) A `fixed_date` block as defined below. This cannot be specified if a `recurrence` block is specified.

* `recurrence` - (Optional) A `recurrence` block as defined below. This cannot be specified if a `fixed_date` block is specified.

---

A `capacity` block supports the following:

* `default` - (Required) The number of instances that are available for scaling if metrics are not available for evaluation. The default is only used if the current instance count is lower than the default. Valid values are between `0` and `1000`.

* `maximum` - (Required) The maximum number of instances for this resource. Valid values are between `0` and `1000`.

-> **NOTE:** The maximum number of instances is also limited by the amount of Cores available in the subscription.

* `minimum` - (Required) The minimum number of instances for this resource. Valid values are between `0` and `1000`.

---

A `rule` block supports the following:

* `metric_trigger` - (Required) A `metric_trigger` block as defined below.

* `scale_action` - (Required) A `scale_action` block as defined below.

---

A `metric_trigger` block supports the following:

* `metric_name` - (Required) The name of the metric that defines what the rule monitors, such as `Percentage CPU` for `Virtual Machine Scale Sets` and `CpuPercentage` for `App Service Plan`.

-> **NOTE:** The allowed value of `metric_name` highly depends on the targeting resource type, please visit [Supported metrics with Azure Monitor](https://docs.microsoft.com/en-us/azure/azure-monitor/platform/metrics-supported) for more details.

* `metric_resource_id` - (Required) The ID of the Resource which the Rule monitors.

* `operator` - (Required) Specifies the operator used to compare the metric data and threshold. Possible values are: `Equals`, `NotEquals`, `GreaterThan`, `GreaterThanOrEqual`, `LessThan`, `LessThanOrEqual`.

* `statistic` - (Required) Specifies how the metrics from multiple instances are combined. Possible values are `Average`, `Min` and `Max`.

* `time_aggregation` - (Required) Specifies how the data that's collected should be combined over time. Possible values include `Average`, `Count`, `Maximum`, `Minimum`, `Last` and `Total`. Defaults to `Average`.

* `time_grain` - (Required) Specifies the granularity of metrics that the rule monitors, which must be one of the pre-defined values returned from the metric definitions for the metric. This value must be between 1 minute and 12 hours an be formatted as an ISO 8601 string.

* `time_window` - (Required) Specifies the time range for which data is collected, which must be greater than the delay in metric collection (which varies from resource to resource). This value must be between 5 minutes and 12 hours and be formatted as an ISO 8601 string.

* `threshold` - (Required) Specifies the threshold of the metric that triggers the scale action.

* `metric_namespace` - (Optional) The namespace of the metric that defines what the rule monitors, such as `microsoft.compute/virtualmachinescalesets` for `Virtual Machine Scale Sets`.

* `dimensions` - (Optional) One or more `dimensions` block as defined below.

---

A `scale_action` block supports the following:

* `cooldown` - (Required) The amount of time to wait since the last scaling action before this action occurs. Must be between 1 minute and 1 week and formatted as a ISO 8601 string.

* `direction` - (Required) The scale direction. Possible values are `Increase` and `Decrease`.

* `type` - (Required) The type of action that should occur. Possible values are `ChangeCount`, `ExactCount` and `PercentChangeCount`.

* `value` - (Required) The number of instances involved in the scaling action. Defaults to `1`.

---

A `fixed_date` block supports the following:

* `end` - (Required) Specifies the end date for the profile, formatted as an RFC3339 date string.

* `start` - (Required) Specifies the start date for the profile, formatted as an RFC3339 date string.

* `timezone` (Optional) The Time Zone of the `start` and `end` times. A list of [possible values can be found here](https://msdn.microsoft.com/en-us/library/azure/dn931928.aspx). Defaults to `UTC`.

---

A `recurrence` block supports the following:

* `timezone` - (Required) The Time Zone used for the `hours` field. A list of [possible values can be found here](https://msdn.microsoft.com/en-us/library/azure/dn931928.aspx). Defaults to `UTC`.

* `days` - (Required) A list of days that this profile takes effect on. Possible values include `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday` and `Sunday`.

* `hours` - (Required) A list containing a single item, which specifies the Hour interval at which this recurrence should be triggered (in 24-hour time). Possible values are from `0` to `23`.

* `minutes` - (Required) A list containing a single item which specifies the Minute interval at which this recurrence should be triggered.

---

A `notification` block supports the following:

* `email` - (Required) A `email` block as defined below.

* `webhook` - (Optional) One or more `webhook` blocks as defined below.

---

A `email` block supports the following:

* `send_to_subscription_administrator` - (Optional) Should email notifications be sent to the subscription administrator? Defaults to `false`.

* `send_to_subscription_co_administrator` - (Optional) Should email notifications be sent to the subscription co-administrator? Defaults to `false`.

* `custom_emails` - (Optional) Specifies a list of custom email addresses to which the email notifications will be sent.

---

A `webhook` block supports the following:

* `service_uri` - (Required) The HTTPS URI which should receive scale notifications.

* `properties` - (Optional) A map of settings.

---

A `dimensions` block supports the following:

* `name` - (Required) The name of the dimension.

* `operator` - (Required) The dimension operator. Possible values are `Equals` and `NotEquals`. `Equals` means being equal to any of the values. `NotEquals` means being not equal to any of the values.

* `values` - (Required) A list of dimension values.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the AutoScale Setting.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the AutoScale Setting.
* `update` - (Defaults to 30 minutes) Used when updating the AutoScale Setting.
* `read` - (Defaults to 5 minutes) Used when retrieving the AutoScale Setting.
* `delete` - (Defaults to 30 minutes) Used when deleting the AutoScale Setting.

## Import

AutoScale Setting can be imported using the `resource id`, e.g.

```
terraform import azurerm_monitor_autoscale_setting.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/microsoft.insights/autoscalesettings/setting1
```
