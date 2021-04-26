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
resource "azurerm_resource_group" "main" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_monitor_action_group" "main" {
  name                = "example-actiongroup"
  resource_group_name = azurerm_resource_group.main.name
  short_name          = "p0action"

  webhook_receiver {
    name        = "callmyapi"
    service_uri = "http://example.com/alert"
  }
}

resource "azurerm_storage_account" "to_monitor" {
  name                     = "examplesa"
  resource_group_name      = azurerm_resource_group.main.name
  location                 = azurerm_resource_group.main.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_monitor_activity_log_alert" "main" {
  name                = "example-activitylogalert"
  resource_group_name = azurerm_resource_group.main.name
  scopes              = [azurerm_resource_group.main.id]
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
* `resource_group_name` - (Required) The name of the resource group in which to create the activity log alert instance.
* `scopes` - (Required) The Scope at which the Activity Log should be applied, for example a the Resource ID of a Subscription or a Resource (such as a Storage Account).
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
* `operation_name` - (Optional) The Resource Manager Role-Based Access Control operation name. Supported operation should be of the form: `<resourceProvider>/<resourceType>/<operation>`.
* `resource_provider` - (Optional) The name of the resource provider monitored by the activity log alert.
* `resource_type` - (Optional) The resource type monitored by the activity log alert.
* `resource_group` - (Optional) The name of resource group monitored by the activity log alert.
* `resource_id` - (Optional) The specific resource monitored by the activity log alert. It should be within one of the `scopes`.
* `caller` - (Optional) The email address or Azure Active Directory identifier of the user who performed the operation.
* `level` - (Optional) The severity level of the event. Possible values are `Verbose`, `Informational`, `Warning`, `Error`, and `Critical`.
* `status` - (Optional) The status of the event. For example, `Started`, `Failed`, or `Succeeded`.
* `sub_status` - (Optional) The sub status of the event.
* `recommendation_type` - (Optional) The recommendation type of the event. It is only allowed when `category` is `Recommendation`.
* `recommendation_category` - (Optional) The recommendation category of the event. Possible values are `Cost`, `Reliability`, `OperationalExcellence` and `Performance`. It is only allowed when `category` is `Recommendation`.
* `recommendation_impact` - (Optional) The recommendation impact of the event. Possible values are `High`, `Medium` and `Low`. It is only allowed when `category` is `Recommendation`.
* `service_health` - (Optional) A block to define fine grain service health settings.

---

A `service_health` block supports the following:

* `events` (Optional) Events this alert will monitor Possible values are `Incident`, `Maintenance`, `Informational`, and `ActionRequired`.
* `locations` (Optional) Locations this alert will monitor. For example, `West Europe`. Defaults to `Global`.
* `services` (Optional) Services this alert will monitor. For example, `Activity Logs & Alerts`, `Action Groups`. Defaults to all Services.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the activity log alert.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Activity Log Alert.
* `update` - (Defaults to 30 minutes) Used when updating the Activity Log Alert.
* `read` - (Defaults to 5 minutes) Used when retrieving the Activity Log Alert.
* `delete` - (Defaults to 30 minutes) Used when deleting the Activity Log Alert.

## Import

Activity log alerts can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_monitor_activity_log_alert.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/microsoft.insights/activityLogAlerts/myalertname
```
