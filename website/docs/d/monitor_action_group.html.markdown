---
subcategory: ""
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_monitor_action_group"
sidebar_current: "docs-azurerm-datasource-monitor-action-group"
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
  value = "${data.azurerm_monitor_action_group.example.id}"
}
```

## Argument Reference

* `name` - (Required) Specifies the name of the Action Group.
* `resource_group_name` - (Required) Specifies the name of the resource group the Action Group is located in.

## Attributes Reference

* `id` - The ID of the Action Group.
* `short_name` - The short name of the action group.
* `enabled` - Whether this action group is enabled.
* `email_receiver` - One or more `email_receiver` blocks as defined below.
* `sms_receiver` - One or more `sms_receiver ` blocks as defined below.
* `webhook_receiver` - One or more `webhook_receiver ` blocks as defined below.

---

`email_receiver` supports the following:

* `name` - The name of the email receiver.
* `email_address` - The email address of this receiver.

---

`sms_receiver` supports the following:

* `name` - The name of the SMS receiver.
* `country_code` - The country code of the SMS receiver.
* `phone_number` - The phone number of the SMS receiver.

---

`webhook_receiver` supports the following:

* `name` - The name of the webhook receiver. 
* `service_uri` - The URI where webhooks should be sent.
