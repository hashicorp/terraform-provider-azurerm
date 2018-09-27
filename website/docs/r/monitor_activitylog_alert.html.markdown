---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_monitor_activitylog_alert"
sidebar_current: "docs-azurerm-resource-monitor-activitylog-alert"
description: |-
  Manages Azure monitor alerts on activity log
---

# azurerm_monitor_activitylog_alert

Manages Azure monitor alerts on activity log.

## Example Usage (Monitor all storage account updates in the current subscription)

```hcl
resource "azurerm_resource_group" "main" {
  name     = "ActivityLogAlertTestRG"
  location = "West US"
}

resource "azurerm_monitor_action_group" "main" {
  name                = "ActivityLogAlertTestAction"
  resource_group_name = "${azurerm_resource_group.main.name}"
  short_name          = "p0action"

  webhook_receiver {
    name        = "callmyapi"
    service_uri = "http://example.com/alert"
  }
}

data "azurerm_client_config" "current" {}

resource "azurerm_monitor_activitylog_alert" "main" {
  name                = "AppServiceStateTestActivityLog"
  resource_group_name = "${azurerm_resource_group.main.name}"
  scopes              = ["/subscriptions/${data.azurerm_client_config.current.subscription_id}"]
  description         = "This alert will monitor all storage account updates in the subscription."

  criteria {
    operation_name = "Microsoft.Storage/storageAccounts/write"
    category       = "Recommendation"
  }

  action {
    action_group_id = "${azurerm_monitor_action_group.main.id}"
  }
}
```

## Example Usage (Monitor one specific storage account updates)

```hcl
resource "azurerm_resource_group" "main" {
  name     = "ActivityLogAlertTestRG"
  location = "West US"
}

resource "azurerm_monitor_action_group" "main" {
  name                = "ActivityLogAlertTestAction"
  resource_group_name = "${azurerm_resource_group.main.name}"
  short_name          = "p0action"

  webhook_receiver {
    name        = "callmyapi"
    service_uri = "http://example.com/alert"
  }
}

resource "azurerm_storage_account" "to_monitor" {
  name                     = "actlogtestmonitoredsa"
  resource_group_name      = "${azurerm_resource_group.main.name}"
  location                 = "${azurerm_resource_group.main.location}"
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_monitor_activitylog_alert" "main" {
  name                = "AppServiceStateTestActivityLog"
  resource_group_name = "${azurerm_resource_group.main.name}"
  scopes              = ["${azurerm_resource_group.main.id}"]
  description         = "This alert will monitor a specific storage account updates."

  criteria {
    resource_id    = "${azurerm_storage_account.to_monitor.id}"
    operation_name = "Microsoft.Storage/storageAccounts/write"
    category       = "Recommendation"
  }

  action {
    action_group_id = "${azurerm_monitor_action_group.main.id}"

    webhook_properties {
      from = "terraform"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the activity log alert. Changing this forces a new resource to be created.
* `resource_group_name` - (Required) The name of the resource group in which to create the activity log alert instance.
* `scopes` - (Required) The individual resource or set of resources for which the alert on activity log is defined. For example, a subscription ID, or a resource group ID.
* `criteria` - (Required) One and only one `criteria` block as defined below. The condition that will cause this alert to activate.
* `action` - (Required) One or more `action` blocks as defined below. The actions that will activate when the condition (defined through `criteria`) is met.
* `enabled` - (Optional) Whether this activity log alert is enabled. Defaults to `true`.
* `description` - (Optional) The description of this activity log alert.
* `tags` - (Optional) A mapping of tags to assign to the resource.

---

`criteria` supports the following:

* `category` - (Required) The category of the operation. Possible values are `Administrative`, `Autoscale`, `Policy`, `Recommendation`, `Security` and `Service Health`.
* `operation_name` - (Required) The Resource Manager Role-Based Access Control operation name. Supported operation should be of the form: `<resourceProvider>/<resourceType>/<operation>`.
* `resource_id` - (Optional) The specific resource monitored by the activity log alert. It should be within one of the `scopes`.
* `caller` - (Optional) The email address or Azure Active Directory identifier of the user who performed the operation.
* `level` - (Optional) The severity level of the event. Possible values are `Verbose`, `Informational`, `Warning`, `Error`, and `Critical`.
* `status` - (Optional) The status of the event. for example, `Started`, `Failed`, or `Succeeded`.
* `sub_status` - (Optional) The sub status of the event.

---

`action` supports the following:

* `action_group_id` - (Required) The resource ID of an action group (`azurerm_monitor_action_group`).
* `webhook_properties` - (Optional) The map of custom string properties to include with the post operation. These data are appended to the webhook payload.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the activity log alert.

## Import

Activity log alerts can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_monitor_activitylog_alert.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/microsoft.insights/activityLogAlerts/myalertname
```
