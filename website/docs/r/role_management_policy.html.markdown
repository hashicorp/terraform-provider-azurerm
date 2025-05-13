---
subcategory: "Authorization"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_role_management_policy"
description: |-
  Manages Azure PIM Role Management Policies.
---

# Resource: azurerm_role_management_policy

Manage a role policy for an Azure Management Group, Subscription, Resource Group or resource.

## Example Usage

### Resource Group

```terraform
resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "East US"
}

data "azurerm_role_definition" "rg_contributor" {
  name  = "Contributor"
  scope = azurerm_resource_group.example.id
}

data "azuread_group" "approvers" {
  display_name = "Example Approver Group"
}

resource "azurerm_role_management_policy" "example" {
  scope              = azurerm_resource_group.test.id
  role_definition_id = data.azurerm_role_definition.contributor.id

  active_assignment_rules {
    expire_after = "P365D"
  }

  eligible_assignment_rules {
    expiration_required = false
  }

  activation_rules {
    maximum_duration = "PT1H"
    require_approval = true
    approval_stage {
      primary_approver {
        object_id = data.azuread_group.approvers.object_id
        type      = "Group"
      }
    }
  }

  notification_rules {
    eligible_assignments {
      approver_notifications {
        notification_level    = "Critical"
        default_recipients    = false
        additional_recipients = ["someone@example.com"]
      }
    }
    eligible_activations {
      assignee_notifications {
        notification_level    = "All"
        default_recipients    = true
        additional_recipients = ["someone.else@example.com"]
      }
    }
  }
}
```

### Management Group

```terraform
resource "azurerm_management_group" "example" {
  name = "example-group"
}

data "azurerm_role_definition" "mg_contributor" {
  name  = "Contributor"
  scope = azurerm_management_group.example.id
}

resource "azurerm_role_management_policy" "example" {
  scope              = azurerm_management_group.example.id
  role_definition_id = data.azurerm_role_definition.mg_contributor.id

  eligible_assignment_rules {
    expiration_required = false
  }

  active_assignment_rules {
    expire_after = "P90D"
  }

  activation_rules {
    maximum_duration = "PT1H"
    require_approval = true
  }

  notification_rules {
    active_assignments {
      admin_notifications {
        notification_level    = "Critical"
        default_recipients    = false
        additional_recipients = ["someone@example.com"]
      }
    }
  }
}
```

## Argument Reference

* `activation_rules` - (Optional) An `activation_rules` block as defined below.
* `active_assignment_rules` - (Optional) An `active_assignment_rules` block as defined below.
* `eligible_assignment_rules` - (Optional) An `eligible_assignment_rules` block as defined below.
* `notification_rules` - (Optional) A `notification_rules` block as defined below.
* `role_definition_id` - (Required) The scoped Role Definition ID of the role for which this policy will apply. Changing this forces a new resource to be created.
* `scope` - (Required) The scope to which this Role Management Policy will apply. Can refer to a management group, a subscription, a resource group or a resource. Changing this forces a new resource to be created.

---

An `activation_rules` block supports the following:

* `approval_stage` - (Optional) An `approval_stage` block as defined below.
* `maximum_duration` - (Optional) The maximum length of time an activated role can be valid, in an ISO8601 Duration format (e.g. `PT8H`). Valid range is `PT30M` to `PT23H30M`, in 30 minute increments, or `PT1D`.
* `require_approval` - (Optional) Is approval required for activation. If `true` an `approval_stage` block must be provided.
* `require_justification` - (Optional) Is a justification required during activation of the role.
* `require_multifactor_authentication` - (Optional) Is multi-factor authentication required to activate the role. Conflicts with `required_conditional_access_authentication_context`.
* `require_ticket_info` - (Optional) Is ticket information requrired during activation of the role.
* `required_conditional_access_authentication_context` - (Optional) The Entra ID Conditional Access context that must be present for activation. Conflicts with `require_multifactor_authentication`.

---

An `active_assignment_rules` block supports the following:

* `expiration_required` - (Optional) Must an assignment have an expiry date. `false` allows permanent assignment.
* `expire_after` - (Optional) The maximum length of time an assignment can be valid, as an ISO8601 duration. Permitted values: `P15D`, `P30D`, `P90D`, `P180D`, or `P365D`.
* `require_justification` - (Optional) Is a justification required to create new assignments.
* `require_multifactor_authentication` - (Optional) Is multi-factor authentication required to create new assignments.
* `require_ticket_info` - (Optional) Is ticket information required to create new assignments.

One of `expiration_required` or `expire_after` must be provided.

---

An `approval_stage` block supports the following:

* One or more `primary_approver` blocks as defined below.

---

An `eligible_assignment_rules` block supports the following:

* `expiration_required`- Must an assignment have an expiry date. `false` allows permanent assignment.
* `expire_after` - The maximum length of time an assignment can be valid, as an ISO8601 duration. Permitted values: `P15D`, `P30D`, `P90D`, `P180D`, or `P365D`.

One of `expiration_required` or `expire_after` must be provided.

---

A `notification_rules` block supports the following:

* `active_assignments` - (Optional) A `notification_target` block as defined below to configure notfications on active role assignments.
* `eligible_activations` - (Optional) A `notification_target` block as defined below for configuring notifications on activation of eligible role.
* `eligible_assignments` - (Optional) A `notification_target` block as defined below to configure notification on eligible role assignments.

At least one `notification_target` block must be provided.

---

A `notification_settings` block supports the following:

* `additional_recipients` - (Optional) A list of additional email addresses that will receive these notifications.
* `default_recipients` - (Required) Should the default recipients receive these notifications.
* `notification_level` - (Required) What level of notifications should be sent. Options are `All` or `Critical`.

---

A `notification_target` block supports the following:

* `admin_notifications` - (Optional) A `notification_settings` block as defined above.
* `approver_notifications` - (Optional) A `notification_settings` block as defined above.
* `assignee_notifications` - (Optional) A `notification_settings` block as defined above.

At least one `notification_settings` block must be provided.

---

A `primary_approver` block supports the following:

* `object_id` - (Required) The ID of the object which will act as an approver.
* `type` - (Required) The type of object acting as an approver. Possible options are `User` and `Group`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` (String) The ID of this policy.
* `name` (String) The name of this policy, which is typically a UUID and may change over time.
* `description` (String) The description of this policy.

## Import

Because these policies are created automatically by Azure, they will auto-import on first use. They can be imported using the `resource id` of the role definition, combined with the scope id, e.g.

```shell
terraform import azurerm_role_management_policy.example "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/roleDefinitions/00000000-0000-0000-0000-000000000000|<scope>"
```

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Role Definition.
* `read` - (Defaults to 5 minutes) Used when retrieving the Role Definition.
* `update` - (Defaults to 30 minutes) Used when updating the Role Definition.
* `delete` - (Defaults to 5 minutes) Used when deleting the Role Definition.
