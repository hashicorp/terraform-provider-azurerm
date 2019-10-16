---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_monitor_action_group"
sidebar_current: "docs-azurerm-resource-monitor-action-group"
description: |-
  Manages an Action Group within Azure Monitor

---

# azurerm_monitor_action_group

Manages an Action Group within Azure Monitor.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "monitoring-resources"
  location = "West US"
}

resource "azurerm_monitor_action_group" "test" {
  name                = "CriticalAlertsAction"
  resource_group_name = "${azurerm_resource_group.test.name}"
  short_name          = "p0action"

  email_receiver {
    name          = "sendtoadmin"
    email_address = "admin@contoso.com"
  }

  email_receiver {
    name          = "sendtodevops"
    email_address = "devops@contoso.com"
  }

  itsm_receiver {         
    name                 = "createorupdateticket"
    workspace_id         = "6eee3a18-aac3-40e4-b98e-1f309f329816"
    connection_id        = "53de6956-42b4-41ba-be3c-b154cdf17b13"
    ticket_configuration = "{}"
    region               = "southcentralus"
  }

  azure_app_push_receiver {
    name          = "pushtoadmin"
    email_address = "admin@contoso.com"
  }

  sms_receiver {
    name         = "oncallmsg"
    country_code = "1"
    phone_number = "1231231234"
  }

  webhook_receiver {
    name                    = "callmyapiaswell"
    service_uri             = "http://example.com/alert"
    use_common_alert_schema = true
  }

  automation_runbook_receiver {   
    name                    = "action_name_1"
    automation_account_id   = "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/rg-runbooks/providers/microsoft.automation/automationaccounts/aaa001"
    runbook_name            = "my runbook"
    webhook_resource_id     = "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/rg-runbooks/providers/microsoft.automation/automationaccounts/aaa001/webhooks/webhook_alert"
    is_global_runbook       = true
    service_uri             = "https://s13events.azure-automation.net/webhooks?token=randomtoken"
  }

  voice_receiver { 
    name         = "remotesupport"
    country_code = "86"
    phone_number = "13888888888"
  }

  logic_app_receiver {
    name = "logicappaction"
    resource_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg-logicapp/providers/Microsoft.Logic/workflows/logicapp"
    callback_url = "https://logicapptriggerurl/..."
  }

  azure_function_receiver {
    name = "funcaction"
    function_app_resource_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg-funcapp/providers/Microsoft.Web/sites/funcapp"
    function_name = "myfunc"
    http_trigger_url = "https://example.com/trigger"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Action Group. Changing this forces a new resource to be created.
* `resource_group_name` - (Required) The name of the resource group in which to create the Action Group instance.
* `short_name` - (Required) The short name of the action group. This will be used in SMS messages.
* `enabled` - (Optional) Whether this action group is enabled. If an action group is not enabled, then none of its receivers will receive communications. Defaults to `true`.
* `email_receiver` - (Optional) One or more `email_receiver` blocks as defined below.
* `itsm_receiver` - (Optional) One or more `itsm_receiver` blocks as defined below.
* `azure_app_push_receiver` - (Optional) One or more `azure_app_push_receiver` blocks as defined below.
* `sms_receiver` - (Optional) One or more `sms_receiver ` blocks as defined below.
* `webhook_receiver` - (Optional) One or more `webhook_receiver ` blocks as defined below.
* `automation_runbook_receiver` - (Optional) One or more `automation_runbook_receiver` blocks as defined below.
* `voice_receiver` - (Optional) One or more `voice_receiver` blocks as defined below.
* `azure_function_receiver` - (Optional) One or more `azure_function_receiver` blocks as defined below.
* `logic_app_receiver` - (Optional) One or more `logic_app_receiver` blocks as defined below.
* `tags` - (Optional) A mapping of tags to assign to the resource.

---

`email_receiver` supports the following:

* `name` - (Required) The name of the email receiver. Names must be unique (case-insensitive) across all receivers within an action group.
* `email_address` - (Required) The email address of this receiver.

---

`itsm_receiver` supports the following:

* `name` - (Required) The name of the ITSM receiver.
* `workspace_id` - (Required) The Azure Log Analytics workspace ID where this connection is defined.
* `connection_id` - (Required) The unique connection identifier of the ITSM connection.
* `ticket_configuration` - (Required) A JSON blob for the configurations of the ITSM action. CreateMultipleWorkItems option will be part of this blob as well.
* `region` - (Required) The region of the workspace.

---

`azure_app_push_receiver` supports the following:

* `name` - (Required) The name of the Azure app push receiver.
* `email_address` - (Required) The email address of the user signed into the mobile app who will receive push notifications from this receiver.

---

`sms_receiver` supports the following:

* `name` - (Required) The name of the SMS receiver. Names must be unique (case-insensitive) across all receivers within an action group.
* `country_code` - (Required) The country code of the SMS receiver.
* `phone_number` - (Required) The phone number of the SMS receiver.

---

`webhook_receiver` supports the following:

* `name` - (Required) The name of the webhook receiver. Names must be unique (case-insensitive) across all receivers within an action group.
* `service_uri` - (Required) The URI where webhooks should be sent.
* `use_common_alert_schema` - (Optional) Enables or disables the common alert schema.

---

`automation_runbook_receiver` supports the following:

* `name` - (Required) The name of the automation runbook receiver. 
* `automation_account_id` - (Required) The automation account ID which holds this runbook and authenticates to Azure resources.
* `runbook_name` - (Required) The name for this runbook.
* `webhook_resource_id` - (Required) The resource id for webhook linked to this runbook.
* `is_global_runbook` - (Required) Indicates whether this instance is global runbook.
* `service_uri` - (Required) The URI where webhooks should be sent.

---

`voice_receiver` supports the following:

* `name` - (Required) The name of the voice receiver. 
* `country_code` - (Required) The country code of the voice receiver.
* `phone_number` - (Required) The phone number of the voice receiver.

---

`azure_function_receiver` supports the following:

* `name` - (Required) The name of the Azure Function receiver.
* `function_app_resouce_id` - (Required) The Azure resource ID of the function app.
* `function_name` - (Required) The function name in the function app.
* `http_trigger_url` - (Required) The http trigger url where http request sent to.

---

`logic_app_receiver` supports the following:

* `name` - (Required) The name of the logic app receiver.
* `resource_id` - (Required) The Azure resource ID of the logic app.
* `callback_url` - (Required) The callback url where http request sent to.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Action Group.

## Import

Action Groups can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_monitor_action_group.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Insights/actionGroups/myagname
```
