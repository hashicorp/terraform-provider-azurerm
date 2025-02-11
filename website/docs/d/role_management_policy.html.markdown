---
subcategory: "Authorization"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_role_management_policy"
description: |-
  Get information about an Azure PIM Role Management Policy.
---

# Data Source: azurerm_role_management_policy

Use this data source to get information on a role policy for an Azure Management Group, Subscription, Resource Group or resource.

## Example Usage

### Resource Group

```terraform
data "azurerm_resource_group" "example" {
  name = "example-rg"
}

data "azurerm_role_definition" "rg_contributor" {
  name  = "Contributor"
  scope = data.azurerm_resource_group.example.id
}

data "azurerm_role_management_policy" "example" {
  scope              = data.azurerm_resource_group.test.id
  role_definition_id = data.azurerm_role_definition.contributor.id
}
```

### Management Group

```terraform
data "azurerm_management_group" "example" {
  name = "example-group"
}

data "azurerm_role_definition" "mg_contributor" {
  name  = "Contributor"
  scope = azurerm_management_group.example.id
}

data "azurerm_role_management_policy" "example" {
  scope              = data.azurerm_management_group.example.id
  role_definition_id = data.azurerm_role_definition.mg_contributor.id
}
```

## Argument Reference

* `role_definition_id` - (Required) The scoped Role Definition ID of the role for which this policy applies.
* `scope` - (Required) The scope to which this Role Management Policy applies. Can refer to a management group, a subscription, a resource group or a resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` (String) The ID of this policy.
* `name` (String) The name of this policy, which is typically a UUID and may change over time.
* `description` (String) The description of this policy.
* `activation_rules` - An `activation_rules` block as defined below.
* `active_assignment_rules` - An `active_assignment_rules` block as defined below.
* `eligible_assignment_rules` - An `eligible_assignment_rules` block as defined below.
* `notification_rules` - A `notification_rules` block as defined below.

---

An `activation_rules` block returns the following:

* `approval_stage` - An `approval_stage` block as defined below.
* `maximum_duration` - (String) The maximum length of time an activated role can be valid, in an ISO8601 Duration format.
* `require_approval` - (Boolean) Is approval required for activation.
* `require_justification` - (Boolean) Is a justification required during activation of the role.
* `require_multifactor_authentication` - (Boolean) Is multi-factor authentication required to activate the role.
* `require_ticket_info` - (Boolean) Is ticket information requrired during activation of the role.
* `required_conditional_access_authentication_context` - (String) The Entra ID Conditional Access context that must be present for activation.

---

An `active_assignment_rules` block returns the following:

* `expiration_required` - (Boolean) Must an assignment have an expiry date.
* `expire_after` - (String) The maximum length of time an assignment can be valid, as an ISO8601 duration.
* `require_justification` - (Boolean) Is a justification required to create new assignments.
* `require_multifactor_authentication` - (Boolean) Is multi-factor authentication required to create new assignments.
* `require_ticket_info` - (Boolean) Is ticket information required to create new assignments.

---

An `approval_stage` block returns the following:

* One or more `primary_approver` blocks as defined below.

---

An `eligible_assignment_rules` block returns the following:

* `expiration_required`- (Boolean) Must an assignment have an expiry date.
* `expire_after` - (String) The maximum length of time an assignment can be valid, as an ISO8601 duration.

---

A `notification_rules` block returns the following:

* `active_assignments` - A `notification_target` block as defined below with the details of notfications on active role assignments.
* `eligible_activations` - A `notification_target` block as defined below with the details of notifications on activation of eligible role.
* `eligible_assignments` - A `notification_target` block as defined below with the details of notifications on eligible role assignments.

---

A `notification_settings` block returns the following:

* `additional_recipients` - A list of additional email addresses that will receive these notifications.
* `default_recipients` - (Boolean) Should the default recipients receive these notifications.
* `notification_level` - (String) What level of notifications should be sent. Either `All` or `Critical`.

---

A `notification_target` block returns the following:

* `admin_notifications` - A `notification_settings` block as defined above.
* `approver_notifications` - A `notification_settings` block as defined above.
* `assignee_notifications` - A `notification_settings` block as defined above.

---

A `primary_approver` block returns the following:

* `object_id` - (String) The ID of the object which will act as an approver.
* `type` - (String) The type of object acting as an approver. Either `User` or `Group`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Role Definition.
