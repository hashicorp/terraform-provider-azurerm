---
subcategory: "Monitor"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_monitor_activity_log_alert"
description: |-
  Manages an Activity Log Alert within Azure Monitor
---

# azurerm_monitor_activity_log_alert

Manages an Activity Log Alert within Azure Monitor.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_monitor_action_group" "main" {
  name                = "example-actiongroup"
  resource_group_name = azurerm_resource_group.example.name
  short_name          = "p0action"

  webhook_receiver {
    name        = "callmyapi"
    service_uri = "http://example.com/alert"
  }
}

resource "azurerm_storage_account" "to_monitor" {
  name                     = "examplesa"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_monitor_activity_log_alert" "main" {
  name                = "example-activitylogalert"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  scopes              = [azurerm_resource_group.example.id]
  description         = "This alert will monitor a specific storage account updates."

  criteria {
    resource_id    = azurerm_storage_account.to_monitor.id
    operation_name = "Microsoft.Storage/storageAccounts/write"
    category       = "Recommendation"
  }

  action {
    action_group_id = azurerm_monitor_action_group.main.id

    webhook_properties = {
      from = "terraform"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the activity log alert. Changing this forces a new resource to be created.
* `resource_group_name` - (Required) The name of the resource group in which to create the activity log alert instance. Changing this forces a new resource to be created.
* `location` - (Required) The Azure Region where the activity log alert rule should exist. Changing this forces a new resource to be created.
* `scopes` - (Required) The Scope at which the Activity Log should be applied. A list of strings which could be a resource group , or a subscription, or a resource ID (such as a Storage Account).
* `criteria` - (Required) A `criteria` block as defined below.
* `action` - (Optional) One or more `action` blocks as defined below.
* `enabled` - (Optional) Should this Activity Log Alert be enabled? Defaults to `true`.
* `description` - (Optional) The description of this activity log alert.
* `tags` - (Optional) A mapping of tags to assign to the resource.

---

An `action` block supports the following:

* `action_group_id` - (Required) The ID of the Action Group can be sourced from [the `azurerm_monitor_action_group` resource](./monitor_action_group.html).
* `webhook_properties` - (Optional) The map of custom string properties to include with the post operation. These data are appended to the webhook payload.

---

A `criteria` block supports the following:

* `category` - (Required) The category of the operation. Possible values are `Administrative`, `Autoscale`, `Policy`, `Recommendation`, `ResourceHealth`, `Security` and `ServiceHealth`.
* `caller` - (Optional) The email address or Azure Active Directory identifier of the user who performed the operation.
* `operation_name` - (Optional) The Resource Manager Role-Based Access Control operation name. Supported operation should be of the form: `<resourceProvider>/<resourceType>/<operation>`.
* `resource_provider` - (Optional) The name of the resource provider monitored by the activity log alert.
* `resource_providers` - (Optional) A list of names of resource providers monitored by the activity log alert.

~> **Note:** `resource_provider` and `resource_providers` are mutually exclusive.

* `resource_type` - (Optional) The resource type monitored by the activity log alert.
* `resource_types` - (Optional) A list of resource types monitored by the activity log alert.

~> **Note:** `resource_type` and `resource_types` are mutually exclusive.

* `resource_group` - (Optional) The name of resource group monitored by the activity log alert.
* `resource_groups` - (Optional) A list of names of resource groups monitored by the activity log alert.

~> **Note:** `resource_group` and `resource_groups` are mutually exclusive.

* `resource_id` - (Optional) The specific resource monitored by the activity log alert. It should be within one of the `scopes`.
* `resource_ids` - (Optional) A list of specific resources monitored by the activity log alert. It should be within one of the `scopes`.

~> **Note:** `resource_id` and `resource_ids` are mutually exclusive.

* `level` - (Optional) The severity level of the event. Possible values are `Verbose`, `Informational`, `Warning`, `Error`, and `Critical`.
* `levels` - (Optional) A list of severity level of the event. Possible values are `Verbose`, `Informational`, `Warning`, `Error`, and `Critical`.

~> **Note:** `level` and `levels` are mutually exclusive.

* `status` - (Optional) The status of the event. For example, `Started`, `Failed`, or `Succeeded`.
* `statuses` - (Optional) A list of status of the event. For example, `Started`, `Failed`, or `Succeeded`.

~> **Note:** `status` and `statuses` are mutually exclusive.

* `sub_status` - (Optional) The sub status of the event.
* `sub_statuses` - (Optional) A list of sub status of the event.

~> **Note:** `sub_status` and `sub_statuses` are mutually exclusive.
 
* `recommendation_type` - (Optional) The recommendation type of the event. It is only allowed when `category` is `Recommendation`.
* `recommendation_category` - (Optional) The recommendation category of the event. Possible values are `Cost`, `Reliability`, `OperationalExcellence`, `HighAvailability` and `Performance`. It is only allowed when `category` is `Recommendation`.
* `recommendation_impact` - (Optional) The recommendation impact of the event. Possible values are `High`, `Medium` and `Low`. It is only allowed when `category` is `Recommendation`.
* `resource_health` - (Optional) A block to define fine grain resource health settings.
* `service_health` - (Optional) A block to define fine grain service health settings.

---

A `resource_health` block supports the following:

* `current` - (Optional) The current resource health statuses that will log an alert. Possible values are `Available`, `Degraded`, `Unavailable` and `Unknown`.
* `previous` - (Optional) The previous resource health statuses that will log an alert. Possible values are `Available`, `Degraded`, `Unavailable` and `Unknown`.
* `reason` - (Optional) The reason that will log an alert. Possible values are `PlatformInitiated` (such as a problem with the resource in an affected region of an Azure incident), `UserInitiated` (such as a shutdown request of a VM) and `Unknown`.

---

A `service_health` block supports the following:

* `events` - (Optional) Events this alert will monitor Possible values are `Incident`, `Maintenance`, `Informational`, `ActionRequired` and `Security`.
* `locations` - (Optional) Locations this alert will monitor. For example, `West Europe`.
* `services` - (Optional) Services this alert will monitor. For example, `Activity Logs & Alerts`, `Action Groups`. Defaults to all Services.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the activity log alert.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Activity Log Alert.
* `read` - (Defaults to 5 minutes) Used when retrieving the Activity Log Alert.
* `update` - (Defaults to 30 minutes) Used when updating the Activity Log Alert.
* `delete` - (Defaults to 30 minutes) Used when deleting the Activity Log Alert.

## Import

Activity log alerts can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_monitor_activity_log_alert.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Insights/activityLogAlerts/myalertname
```
