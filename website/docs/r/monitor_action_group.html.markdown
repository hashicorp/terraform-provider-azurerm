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
  location = "West US"
}

resource "azuread_application" "example" {
  name            = "example-app"
  identifier_uris = ["https://uri"]
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
    automation_account_id   = "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/rg-runbooks/providers/microsoft.automation/automationaccounts/aaa001"
    runbook_name            = "my runbook"
    webhook_resource_id     = "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/rg-runbooks/providers/microsoft.automation/automationaccounts/aaa001/webhooks/webhook_alert"
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

  itsm_receiver {
    name                 = "createorupdateticket"
    workspace_id         = "6eee3a18-aac3-40e4-b98e-1f309f329816"
    connection_id        = "53de6956-42b4-41ba-be3c-b154cdf17b13"
    ticket_configuration = "{}"
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
* `resource_group_name` - (Required) The name of the resource group in which to create the Action Group instance.
* `short_name` - (Required) The short name of the action group. This will be used in SMS messages.
* `enabled` - (Optional) Whether this action group is enabled. If an action group is not enabled, then none of its receivers will receive communications. Defaults to `true`.
* `arm_role_receiver` - (Optional) One or more `arm_role_receiver` blocks as defined below.
* `automation_runbook_receiver` - (Optional) One or more `automation_runbook_receiver` blocks as defined below.
* `azure_app_push_receiver` - (Optional) One or more `azure_app_push_receiver` blocks as defined below.
* `azure_function_receiver` - (Optional) One or more `azure_function_receiver` blocks as defined below.
* `email_receiver` - (Optional) One or more `email_receiver` blocks as defined below.
* `itsm_receiver` - (Optional) One or more `itsm_receiver` blocks as defined below.
* `logic_app_receiver` - (Optional) One or more `logic_app_receiver` blocks as defined below.
* `sms_receiver` - (Optional) One or more `sms_receiver` blocks as defined below.
* `voice_receiver` - (Optional) One or more `voice_receiver` blocks as defined below.
* `webhook_receiver` - (Optional) One or more `webhook_receiver` blocks as defined below.
* `tags` - (Optional) A mapping of tags to assign to the resource.

---

`arm_role_receiver` supports the following:

* `name` - (Required) The name of the ARM role receiver.
* `role_id` - (Required) The arm role id.
* `use_common_alert_schema` - (Optional) Enables or disables the common alert schema.

---

`automation_runbook_receiver` supports the following:

* `name` - (Required) The name of the automation runbook receiver.
* `automation_account_id` - (Required) The automation account ID which holds this runbook and authenticates to Azure resources.
* `runbook_name` - (Required) The name for this runbook.
* `webhook_resource_id` - (Required) The resource id for webhook linked to this runbook.
* `is_global_runbook` - (Required) Indicates whether this instance is global runbook.
* `service_uri` - (Required) The URI where webhooks should be sent.
* `use_common_alert_schema` - (Optional) Enables or disables the common alert schema.

---

`azure_app_push_receiver` supports the following:

* `name` - (Required) The name of the Azure app push receiver.
* `email_address` - (Required) The email address of the user signed into the mobile app who will receive push notifications from this receiver.

---

`azure_function_receiver` supports the following:

* `name` - (Required) The name of the Azure Function receiver.
* `function_app_resource_id` - (Required) The Azure resource ID of the function app.
* `function_name` - (Required) The function name in the function app.
* `http_trigger_url` - (Required) The http trigger url where http request sent to.
* `use_common_alert_schema` - (Optional) Enables or disables the common alert schema.

---

`email_receiver` supports the following:

* `name` - (Required) The name of the email receiver. Names must be unique (case-insensitive) across all receivers within an action group.
* `email_address` - (Required) The email address of this receiver.
* `use_common_alert_schema` - (Optional) Enables or disables the common alert schema.

---

`itsm_receiver` supports the following:

* `name` - (Required) The name of the ITSM receiver.
* `workspace_id` - (Required) The Azure Log Analytics workspace ID where this connection is defined.
* `connection_id` - (Required) The unique connection identifier of the ITSM connection.
* `ticket_configuration` - (Required) A JSON blob for the configurations of the ITSM action. CreateMultipleWorkItems option will be part of this blob as well.
* `region` - (Required) The region of the workspace.

---

`logic_app_receiver` supports the following:

* `name` - (Required) The name of the logic app receiver.
* `resource_id` - (Required) The Azure resource ID of the logic app.
* `callback_url` - (Required) The callback url where http request sent to.
* `use_common_alert_schema` - (Optional) Enables or disables the common alert schema.

---

`sms_receiver` supports the following:

* `name` - (Required) The name of the SMS receiver. Names must be unique (case-insensitive) across all receivers within an action group.
* `country_code` - (Required) The country code of the SMS receiver.
* `phone_number` - (Required) The phone number of the SMS receiver.

---

`voice_receiver` supports the following:

* `name` - (Required) The name of the voice receiver.
* `country_code` - (Required) The country code of the voice receiver.
* `phone_number` - (Required) The phone number of the voice receiver.

---

`webhook_receiver` supports the following:

* `name` - (Required) The name of the webhook receiver. Names must be unique (case-insensitive) across all receivers within an action group.
* `service_uri` - (Required) The URI where webhooks should be sent.
* `use_common_alert_schema` - (Optional) Enables or disables the common alert schema.
* `use_aad_auth` - (Optional) Use AAD authentication?
~> **NOTE:** Before adding a secure webhook receiver by enabling `use_aad_auth` and setting `aad_auth_object_id`, please read [the configuration instruction of the AAD application](https://docs.microsoft.com/en-us/azure/azure-monitor/platform/action-groups#secure-webhook). 
* `aad_auth_object_id` - (Optional) The webhook app object Id for aad auth. Required when `use_aad_auth` is `true`.
* `aad_auth_identifier_uri` - (Optional) The identifier uri for aad auth.
* `aad_auth_tenant_id` - (Optional) The tenant id for aad auth.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Action Group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Action Group.
* `update` - (Defaults to 30 minutes) Used when updating the Action Group.
* `read` - (Defaults to 5 minutes) Used when retrieving the Action Group.
* `delete` - (Defaults to 30 minutes) Used when deleting the Action Group.

## Import

Action Groups can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_monitor_action_group.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Insights/actionGroups/myagname
```
