---
subcategory: "Monitor"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_monitor_action_group"
description: |-
  Get information about the specified Action Group.
---

# Data Source: azurerm_monitor_action_group

Use this data source to access the properties of an Action Group.

## Example Usage

```hcl
data "azurerm_monitor_action_group" "example" {
  resource_group_name = "terraform-example-rg"
  name                = "tfex-actiongroup"
}

output "action_group_id" {
  value = data.azurerm_monitor_action_group.example.id
}
```

## Argument Reference

* `name` - Specifies the name of the Action Group.
* `resource_group_name` - Specifies the name of the resource group the Action Group is located in.

## Attributes Reference

* `id` - The ID of the Action Group.
* `short_name` - The short name of the action group.
* `enabled` - Whether this action group is enabled.
* `arm_role_receiver` - One or more `arm_role_receiver` blocks as defined below.
* `automation_runbook_receiver` - One or more `automation_runbook_receiver` blocks as defined below.
* `azure_app_push_receiver` - One or more `azure_app_push_receiver` blocks as defined below.
* `azure_function_receiver` - One or more `azure_function_receiver` blocks as defined below.
* `email_receiver` - One or more `email_receiver` blocks as defined below.
* `itsm_receiver` - One or more `itsm_receiver` blocks as defined below.
* `logic_app_receiver` - One or more `logic_app_receiver` blocks as defined below.
* `sms_receiver` - One or more `sms_receiver` blocks as defined below.
* `webhook_receiver` - One or more `webhook_receiver` blocks as defined below.
* `voice_receiver` - One or more `voice_receiver` blocks as defined below.

---

`arm_role_receiver` supports the following:

* `name` - The name of the ARM role receiver.
* `role_id` - The arm role id.
* `use_common_alert_schema` - Indicates whether to use common alert schema.

---

`automation_runbook_receiver` supports the following:

* `name` - The name of the automation runbook receiver.
* `automation_account_id` - The automation account ID which holds this runbook and authenticates to Azure resources.
* `runbook_name` - The name for this runbook.
* `webhook_resource_id` - The resource id for webhook linked to this runbook.
* `is_global_runbook` - Indicates whether this instance is global runbook.
* `service_uri` - The URI where webhooks should be sent.
* `use_common_alert_schema` - Indicates whether to use common alert schema.

---

`azure_app_push_receiver` supports the following:

* `name` - The name of the Azure app push receiver.
* `email_address` - The email address of the user signed into the mobile app who will receive push notifications from this receiver.

---

`azure_function_receiver` supports the following:

* `name` - The name of the Azure Function receiver.
* `function_app_resource_id` - The Azure resource ID of the function app.
* `function_name` - The function name in the function app.
* `http_trigger_url` - The http trigger url where http request sent to.
* `use_common_alert_schema` - Indicates whether to use common alert schema.

---

`email_receiver` supports the following:

* `name` - The name of the email receiver.
* `email_address` - The email address of this receiver.
* `use_common_alert_schema` - Indicates whether to use common alert schema.

---

`itsm_receiver` supports the following:

* `name` - The name of the ITSM receiver.
* `workspace_id` - The Azure Log Analytics workspace ID where this connection is defined.
* `connection_id` - The unique connection identifier of the ITSM connection.
* `ticket_configuration` - A JSON blob for the configurations of the ITSM action. CreateMultipleWorkItems option will be part of this blob as well.
* `region` - The region of the workspace.

---

`logic_app_receiver` supports the following:

* `name` - The name of the logic app receiver.
* `resource_id` - The Azure resource ID of the logic app.
* `callback_url` - The callback url where http request sent to.
* `use_common_alert_schema` - Indicates whether to use common alert schema.

---

`sms_receiver` supports the following:

* `name` - The name of the SMS receiver.
* `country_code` - The country code of the SMS receiver.
* `phone_number` - The phone number of the SMS receiver.

---

`voice_receiver` supports the following:

* `name` - The name of the voice receiver.
* `country_code` - The country code of the voice receiver.
* `phone_number` - The phone number of the voice receiver.

---

`webhook_receiver` supports the following:

* `name` - The name of the webhook receiver.
* `service_uri` - The URI where webhooks should be sent.
* `use_common_alert_schema` - Indicates whether to use common alert schema.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Action Group.
