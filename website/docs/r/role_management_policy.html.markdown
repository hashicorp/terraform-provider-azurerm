---
subcategory: "Authorization"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_role_management_policy"
description: |-
  Manages a Role Management Policy.
---

# azurerm_role_management_policy

Manages a Role Management Policy.

## Example Usage

```hcl
data "azurerm_subscription" "primary" {}

data "azurerm_role_definition" "test" {
  name = "Monitoring Reader"
}

resource "azurerm_role_management_policy" "test" {
  scope              = data.azurerm_subscription.primary.id
  role_definition_id = "${data.azurerm_subscription.primary.id}${data.azurerm_role_definition.test.id}"

  activation {
    maximum_duration_hours = 12
  }
}
```

## Arguments Reference

The following arguments are supported:

* `role_definition_id` - (Required) The ID of the role definition id.

* `scope` - (Required) The scope.

---

* `activation` - (Optional) A `activation` block as defined below.

* `assignment` - (Optional) A `assignment` block as defined below.

* `notifications` - (Optional) A `notifications` block as defined below.

---

A `activation` block supports the following:

* `approvers` - (Optional) A `approvers` block as defined below.

* `maximum_duration_hours` - (Optional) The maximum duration in hours for an activation..

* `require_justification` - (Optional) Is Justification required for an activation?.

* `require_multi_factor_authentication` - (Optional) Is Multi Factor Authentication required for an activation?

* `require_ticket_information` - (Optional) Is Ticket Information required for an activation?.

---

A `active` block supports the following:

* `allow_permanent` - (Optional) Allow permanent active assignment. Conflicts with `assignment.0.active.0.expire_after_days`,`assignment.0.active.0.expire_after_hours`

* `expire_after_days` - (Optional) The number of days after an active assignments is expired. Conflicts with `assignment.0.active.0.allow_permanent`,`assignment.0.active.0.expire_after_hours`

* `expire_after_hours` - (Optional) The number of hours after an active assignments is expired. Conflicts with `assignment.0.active.0.allow_permanent`,`assignment.0.active.0.expire_after_days`

* `require_justification` - (Optional) Is Justification required for an active assignment?

* `require_multi_factor_authentication` - (Optional) Is Multi Factor Authentication required for an active assignment?

---

A `approvers` block supports the following:

* `group` - (Optional) One or more `group` blocks as defined below.

* `user` - (Optional) One or more `user` blocks as defined below.

---

A `assigned_user` block supports the following:

* `additional_recipients` - (Optional) List of additional recipients to email notifications.

* `critical_emails_only` - (Optional) Will critical emails only be sent?

* `default_recipients` - (Optional) Will notifications be sent to the default recipients?

---

A `assignment` block supports the following:

* `active` - (Optional) A `active` block as defined above.

* `eligible` - (Optional) A `eligible` block as defined below.

---

A `eligible` block supports the following:

* `allow_permanent` - (Optional) Allow permanent eligible assignment. Conflicts with `assignment.0.eligible.0.expire_after_hours`,`assignment.0.eligible.0.expire_after_days`

* `expire_after_days` - (Optional) The number of days after an eligible assignments is expired. Conflicts with `assignment.0.eligible.0.allow_permanent`,`assignment.0.eligible.0.expire_after_hours`

* `expire_after_hours` - (Optional) The number of hours after an eligible assignments is expired. Conflicts with `assignment.0.eligible.0.allow_permanent`,`assignment.0.eligible.0.expire_after_days`

---

A `eligible_member_activate` block supports the following:

* `assigned_user` - (Optional) A `assigned_user` block as defined above.

* `request_for_extension_or_approval` - (Optional) A `request_for_extension_or_approval` block as defined below.

* `role_assignment_alert` - (Optional) A `role_assignment_alert` block as defined below.

---

A `group` block supports the following:

* `id` - (Required) The object id of a group.

* `name` - (Required) The name of the group.

---

A `member_assigned_active` block supports the following:

* `assigned_user` - (Optional) A `assigned_user` block as defined above.

* `request_for_extension_or_approval` - (Optional) A `request_for_extension_or_approval` block as defined below.

* `role_assignment_alert` - (Optional) A `role_assignment_alert` block as defined below.

---

A `member_assigned_eligible` block supports the following:

* `assigned_user` - (Optional) A `assigned_user` block as defined above.

* `request_for_extension_or_approval` - (Optional) A `request_for_extension_or_approval` block as defined below.

* `role_assignment_alert` - (Optional) A `role_assignment_alert` block as defined below.

---

A `notifications` block supports the following:

* `eligible_member_activate` - (Optional) A `eligible_member_activate` block as defined above.

* `member_assigned_active` - (Optional) A `member_assigned_active` block as defined above.

* `member_assigned_eligible` - (Optional) A `member_assigned_eligible` block as defined above.

---

A `request_for_extension_or_approval` block supports the following:

* `additional_recipients` - (Optional) List of additional recipients to email notifications.

* `critical_emails_only` - (Optional) Will critical emails only be sent?

* `default_recipients` - (Optional) Will notifications be sent to the default recipients?

---

A `role_assignment_alert` block supports the following:

* `additional_recipients` - (Optional) List of additional recipients to email notifications.

* `critical_emails_only` - (Optional) Will critical emails only be sent?

* `default_recipients` - (Optional) Will notifications be sent to the default recipients?

---

A `user` block supports the following:

* `id` - (Required) The object id of a user.

* `name` - (Required) The name of the user.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Role Management Policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Role Management Policy.
* `read` - (Defaults to 5 minutes) Used when retrieving the Role Management Policy.
* `update` - (Defaults to 30 minutes) Used when updating the Role Management Policy.
* `delete` - (Defaults to 30 minutes) Used when deleting the Role Management Policy.

## Import

Role Management Policies can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_role_management_policy.example /subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/roleManagementPolicies/00000000-0000-0000-0000-000000000000|/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/roleDefinitions/00000000-0000-0000-0000-000000000000
```

-> **NOTE:** This ID is specific to Terraform - and is of the format `{roleManagementPolicyId}|{roleDefinitionId}`.
