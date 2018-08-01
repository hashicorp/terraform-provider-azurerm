---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_action_group"
sidebar_current: "docs-azurerm-resource-action-group"
description: |-
  Manages an Action Group of Azure monitoring service

---

# azurerm_action_group

Manages an Action Group of Azure monitoring service.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "monitoring-resources"
  location = "West US"
}

resource "azurerm_action_group" "test" {
  name                = "CriticalAlertsAction"
  location            = "Global"
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

  sms_receiver {
    name         = "oncallmsg"
    country_code = "1"
    phone_number = "1231231234"
  }

  webhook_receiver {
    name        = "callmyapiaswell"
    service_uri = "http://example.com/alert"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Action Group. Changing this forces a new resource to be created.
* `location` - (Required) The location of this Action Group. The only possible value is `Global`.
* `resource_group_name` - (Required) The name of the resource group in which to create the Action Group instance.
* `short_name` - (Required) The short name of the action group. This will be used in SMS messages.
* `enabled` - (Optional) Whether this action group is enabled. If an action group is not enabled, then none of its receivers will receive communications. Defaults to `true`.
* `email_receiver` - (Optional) The list of `email_receiver` blocks as defined below that are part of this action group.
* `sms_receiver` - (Optional) The list of `sms_receiver` blocks as defined below that are part of this action group.
* `webhook_receiver` - (Optional) The list of `webhook_receiver` blocks as defined below that are part of this action group.
* `tags` - (Optional) A mapping of tags to assign to the resource.

---

`email_receiver` supports the following:

* `name` - (Required) The name of the email receiver. Names must be unique (case-insensitive) across all receivers within an action group.
* `email_address` - (Required) The email address of this receiver.

---

`sms_receiver` supports the following:

* `name` - (Required) The name of the SMS receiver. Names must be unique (case-insensitive) across all receivers within an action group.
* `country_code` - (Required) The country code of the SMS receiver.
* `phone_number` - (Required) The phone number of the SMS receiver.

---

`webhook_receiver` supports the following:

* `name` - (Required) The name of the webhook receiver. Names must be unique (case-insensitive) across all receivers within an action group.
* `service_uri` - (Required) The URI where webhooks should be sent.

## Attributes Reference

The following attributes are exported:

* `id` - The Route ID.

## Import

Action Groups can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_action_group.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Insights/actionGroups/myagname
```
