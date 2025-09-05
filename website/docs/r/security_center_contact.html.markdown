---
subcategory: "Security Center"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_security_center_contact"
description: |-
  Manages the subscription's Security Center Contact.
---

# azurerm_security_center_contact

Manages the subscription's Security Center Contact.

~> **Note:** Owner access permission is required.

## Example Usage

```hcl
resource "azurerm_security_center_contact" "example" {
  name  = "contact"
  email = "contact@example.com"
  phone = "+1-555-555-5555"

  alert_notifications = true
  alerts_to_admins    = true
}
```

## Arguments Reference

The following arguments are supported:

* `alert_notifications` - (Required) Whether to send security alerts notifications to the security contact.

* `alerts_to_admins` - (Required) Whether to send security alerts notifications to subscription admins.

* `email` - (Required) The email of the Security Center Contact.

* `name` - (Required) The name of the Security Center Contact. Changing this forces a new Security Center Contact to be created.

---

* `phone` - (Optional) The phone number of the Security Center Contact.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Security Center Contact.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the Security Center Contact.
* `read` - (Defaults to 5 minutes) Used when retrieving the Security Center Contact.
* `update` - (Defaults to 1 hour) Used when updating the Security Center Contact.
* `delete` - (Defaults to 1 hour) Used when deleting the Security Center Contact.

## Import

Security Center Contacts can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_security_center_contact.example /subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Security/securityContacts/default1
```
