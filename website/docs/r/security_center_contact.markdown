---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_security_center_contact"
sidebar_current: "docs-azurerm-security-center-contact"
description: |-
    Manages the subscription's Security Center Contact.
---

# azurerm_security_center_contact

Manages the subscription's Security Center Contact.

~> **NOTE:** Owner access permission is required. 

## Example Usage

```hcl
resource "azurerm_security_center_contact" "example" {
  email = "contact@example.com"
  phone = "+1-555-555-5555"

  alert_notifications = true
  alerts_to_admins    = true
}
```

## Argument Reference

The following arguments are supported:

* `email` - (Required) The email of the Security Center Contact.
* `phone` - (Optional) The phone number of the Security Center Contact.
* `alert_notifications` - (Required) Whether to send security alerts notifications to the security contact.
* `alerts_to_admins` - (Required) Whether to send security alerts notifications to subscription admins.

## Attributes Reference

The following attributes are exported:

* `id` - The Security Center Contact ID.

## Import

The contact can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_security_center_contact.example /subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Security/securityContacts/default1
```
