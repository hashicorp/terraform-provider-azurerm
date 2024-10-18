---
subcategory: "Monitor"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_monitor_action_group"
description: |-
  Manages an Action Group within Azure Monitor

---

# azurerm_monitor_action_group

Manages an Action Group within Azure Monitor.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "monitoring-resources"
  location = "West Europe"
}

data "azurerm_client_config" "current" {
}

resource "azurerm_log_analytics_workspace" "example" {
  name                = "workspace-01"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_monitor_action_group" "example" {
  name                = "CriticalAlertsAction"
  resource_group_name = azurerm_resource_group.example.name
  short_name          = "p0action"

  arm_role_receiver {
    name                    = "armroleaction"
    role_id                 = "de139f84-1756-47ae-9be6-808fbbe84772"
    use_common_alert_schema = true
  }

  automation_runbook_receiver {
    name                    = "action_name_1"
    automation_account_id   = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg-runbooks/providers/Microsoft.Automation/automationAccounts/aaa001"
    runbook_name            = "my runbook"
    webhook_resource_id     = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg-runbooks/providers/Microsoft.Automation/automationAccounts/aaa001/webHooks/webhook_alert"
    is_global_runbook       = true
    service_uri             = "https://s13events.azure-automation.net/webhooks?token=randomtoken"
    use_common_alert_schema = true
  }

  azure_app_push_receiver {
    name          = "pushtoadmin"
    email_address = "admin@contoso.com"
  }

  azure_function_receiver {
    name                     = "funcaction"
    function_app_resource_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg-funcapp/providers/Microsoft.Web/sites/funcapp"
    function_name            = "myfunc"
    http_trigger_url         = "https://example.com/trigger"
    use_common_alert_schema  = true
  }

  email_receiver {
    name          = "sendtoadmin"
    email_address = "admin@contoso.com"
  }

  email_receiver {
    name                    = "sendtodevops"
    email_address           = "devops@contoso.com"
    use_common_alert_schema = true
  }

  event_hub_receiver {
    name                    = "sendtoeventhub"
    event_hub_namespace     = "eventhubnamespace"
    event_hub_name          = "eventhub1"
    subscription_id         = "00000000-0000-0000-0000-000000000000"
    use_common_alert_schema = false
  }

  itsm_receiver {
    name                 = "createorupdateticket"
    workspace_id         = "${data.azurerm_client_config.current.subscription_id}|${azurerm_log_analytics_workspace.example.workspace_id}"
    connection_id        = "53de6956-42b4-41ba-be3c-b154cdf17b13"
    ticket_configuration = "{\"PayloadRevision\":0,\"WorkItemType\":\"Incident\",\"UseTemplate\":false,\"WorkItemData\":\"{}\",\"CreateOneWIPerCI\":false}"
    region               = "southcentralus"
  }

  logic_app_receiver {
    name                    = "logicappaction"
    resource_id             = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg-logicapp/providers/Microsoft.Logic/workflows/logicapp"
    callback_url            = "https://logicapptriggerurl/..."
    use_common_alert_schema = true
  }

  sms_receiver {
    name         = "oncallmsg"
    country_code = "1"
    phone_number = "1231231234"
  }

  voice_receiver {
    name         = "remotesupport"
    country_code = "86"
    phone_number = "13888888888"
  }

  webhook_receiver {
    name                    = "callmyapiaswell"
    service_uri             = "http://example.com/alert"
    use_common_alert_schema = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Action Group. Changing this forces a new resource to be created. 
* `resource_group_name` - (Required) The name of the resource group in which to create the Action Group instance. Changing this forces a new resource to be created.
* `short_name` - (Required) The short name of the action group. This will be used in SMS messages.
* `enabled` - (Optional) Whether this action group is enabled. If an action group is not enabled, then none of its receivers will receive communications. Defaults to `true`.
* `arm_role_receiver` - (Optional) One or more `arm_role_receiver` blocks as defined below.
* `automation_runbook_receiver` - (Optional) One or more `automation_runbook_receiver` blocks as defined below.
* `azure_app_push_receiver` - (Optional) One or more `azure_app_push_receiver` blocks as defined below.
* `azure_function_receiver` - (Optional) One or more `azure_function_receiver` blocks as defined below.
* `email_receiver` - (Optional) One or more `email_receiver` blocks as defined below.
* `event_hub_receiver` - (Optional) One or more `event_hub_receiver` blocks as defined below.
* `itsm_receiver` - (Optional) One or more `itsm_receiver` blocks as defined below.
* `location` - (Optional) The Azure Region where the Action Group should exist. Changing this forces a new Action Group to be created. Defaults to `global`.
* `logic_app_receiver` - (Optional) One or more `logic_app_receiver` blocks as defined below.
* `sms_receiver` - (Optional) One or more `sms_receiver` blocks as defined below.
* `voice_receiver` - (Optional) One or more `voice_receiver` blocks as defined below.
* `webhook_receiver` - (Optional) One or more `webhook_receiver` blocks as defined below.
* `tags` - (Optional) A mapping of tags to assign to the resource.

---

The `arm_role_receiver` block supports the following:

* `name` - (Required) The name of the ARM role receiver.
* `role_id` - (Required) The arm role id.
* `use_common_alert_schema` - (Optional) Enables or disables the common alert schema.

---

The `automation_runbook_receiver` block supports the following:

* `name` - (Required) The name of the automation runbook receiver.
* `automation_account_id` - (Required) The automation account ID which holds this runbook and authenticates to Azure resources.
* `runbook_name` - (Required) The name for this runbook.
* `webhook_resource_id` - (Required) The resource id for webhook linked to this runbook.
* `is_global_runbook` - (Required) Indicates whether this instance is global runbook.
* `service_uri` - (Required) The URI where webhooks should be sent.
* `use_common_alert_schema` - (Optional) Enables or disables the common alert schema.

---

The `azure_app_push_receiver` block supports the following:

* `name` - (Required) The name of the Azure app push receiver.
* `email_address` - (Required) The email address of the user signed into the mobile app who will receive push notifications from this receiver.

---

The `azure_function_receiver` block supports the following:

* `name` - (Required) The name of the Azure Function receiver.
* `function_app_resource_id` - (Required) The Azure resource ID of the function app.
* `function_name` - (Required) The function name in the function app.
* `http_trigger_url` - (Required) The HTTP trigger url where HTTP request sent to.
* `use_common_alert_schema` - (Optional) Enables or disables the common alert schema.

---

The `email_receiver` block supports the following:

* `name` - (Required) The name of the email receiver. Names must be unique (case-insensitive) across all receivers within an action group.
* `email_address` - (Required) The email address of this receiver.
* `use_common_alert_schema` - (Optional) Enables or disables the common alert schema.

---

The `event_hub_receiver` block supports the following:

* `name` - (Required) The name of the EventHub Receiver, must be unique within action group.
* `event_hub_name` - (Optional) The name of the specific Event Hub queue.
* `event_hub_namespace` - (Optional) The namespace name of the Event Hub.
* `subscription_id` - (Optional) The ID for the subscription containing this Event Hub. Default to the subscription ID of the Action Group.
* `tenant_id` - (Optional) The Tenant ID for the subscription containing this Event Hub.
* `use_common_alert_schema` - (Optional) Indicates whether to use common alert schema.

---

The `itsm_receiver` block supports the following:

* `name` - (Required) The name of the ITSM receiver.
* `workspace_id` - (Required) The Azure Log Analytics workspace ID where this connection is defined. Format is `<subscription id>|<workspace id>`, for example `00000000-0000-0000-0000-000000000000|00000000-0000-0000-0000-000000000000`.
* `connection_id` - (Required) The unique connection identifier of the ITSM connection.
* `ticket_configuration` - (Required) A JSON blob for the configurations of the ITSM action. CreateMultipleWorkItems option will be part of this blob as well.
* `region` - (Required) The region of the workspace.

-> **NOTE** `ticket_configuration` should be JSON blob with `PayloadRevision` and `WorkItemType` keys (e.g., `ticket_configuration="{\"PayloadRevision\":0,\"WorkItemType\":\"Incident\"}"`), and `ticket_configuration="{}"` will return an error, see more at this [REST API issue](https://github.com/Azure/azure-rest-api-specs/issues/20488)

---

The `logic_app_receiver` block supports the following:

* `name` - (Required) The name of the logic app receiver.
* `resource_id` - (Required) The Azure resource ID of the logic app.
* `callback_url` - (Required) The callback url where HTTP request sent to.
* `use_common_alert_schema` - (Optional) Enables or disables the common alert schema.

---

The `sms_receiver` block supports the following:

* `name` - (Required) The name of the SMS receiver. Names must be unique (case-insensitive) across all receivers within an action group.
* `country_code` - (Required) The country code of the SMS receiver.
* `phone_number` - (Required) The phone number of the SMS receiver.

---

The `voice_receiver` block supports the following:

* `name` - (Required) The name of the voice receiver.
* `country_code` - (Required) The country code of the voice receiver.
* `phone_number` - (Required) The phone number of the voice receiver.

---

The `webhook_receiver` block supports the following:

* `name` - (Required) The name of the webhook receiver. Names must be unique (case-insensitive) across all receivers within an action group.
* `service_uri` - (Required) The URI where webhooks should be sent.
* `use_common_alert_schema` - (Optional) Enables or disables the common alert schema.
* `aad_auth` - (Optional) The `aad_auth` block as defined below.

~> **NOTE:** Before adding a secure webhook receiver by setting `aad_auth`, please read [the configuration instruction of the AAD application](https://docs.microsoft.com/azure/azure-monitor/platform/action-groups#secure-webhook).

---

The `aad_auth` block supports the following:.

* `object_id` - (Required) The webhook application object Id for AAD auth.
* `identifier_uri` - (Optional) The identifier URI for AAD auth.
* `tenant_id` - (Optional) The tenant id for AAD auth.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Action Group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Action Group.
* `update` - (Defaults to 30 minutes) Used when updating the Action Group.
* `read` - (Defaults to 5 minutes) Used when retrieving the Action Group.
* `delete` - (Defaults to 30 minutes) Used when deleting the Action Group.

## Import

Action Groups can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_monitor_action_group.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Insights/actionGroups/myagname
```
