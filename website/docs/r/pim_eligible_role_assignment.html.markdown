---
subcategory: "Authorization"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_pim_eligible_role_assignment"
description: |-
  Manages a PIM Eligible Role Assignment.
---

# azurerm_pim_eligible_role_assignment

Manages a PIM Eligible Role Assignment.

## Example Usage (Subscription)

```hcl
data "azurerm_subscription" "primary" {}

data "azurerm_client_config" "example" {}

data "azurerm_role_definition" "example" {
  name = "Reader"
}

resource "time_static" "example" {}

resource "azurerm_pim_eligible_role_assignment" "example" {
  scope              = data.azurerm_subscription.primary.id
  role_definition_id = "${data.azurerm_subscription.primary.id}${data.azurerm_role_definition.example.id}"
  principal_id       = data.azurerm_client_config.example.object_id

  schedule {
    start_date_time = time_static.example.rfc3339
    expiration {
      duration_hours = 8
    }
  }

  justification = "Expiration Duration Set"

  ticket {
    number = "1"
    system = "example ticket system"
  }
}
```

## Example Usage (Management Group)

```hcl
data "azurerm_client_config" "example" {}

data "azurerm_role_definition" "example" {
  name = "Reader"
}

resource "azurerm_management_group" "example" {
  name = "Example-Management-Group"
}

resource "time_static" "example" {}

resource "azurerm_pim_eligible_role_assignment" "example" {
  scope              = azurerm_management_group.example.id
  role_definition_id = data.azurerm_role_definition.example.id
  principal_id       = data.azurerm_client_config.example.object_id

  schedule {
    start_date_time = time_static.example.rfc3339
    expiration {
      duration_hours = 8
    }
  }

  justification = "Expiration Duration Set"

  ticket {
    number = "1"
    system = "example ticket system"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `principal_id` - (Required) Object ID of the principal for this eligible role assignment. Changing this forces a new resource to be created.

* `role_definition_id` - (Required) The role definition ID for this eligible role assignment. Changing this forces a new resource to be created.

* `scope` - (Required) The scope for this eligible role assignment, should be a valid resource ID. Changing this forces a new resource to be created.

---

* `justification` - (Optional) The justification of the role assignment. Changing this forces a new resource to be created.

* `schedule` - (Optional) A `schedule` block as defined below. Changing this forces a new resource to be created.

* `ticket` - (Optional) A `ticket` block as defined below. Changing this forces a new resource to be created.


* `condition` - (Optional) The condition that limits the resources that the role can be assigned to. See the [official conditions documentation](https://learn.microsoft.com/en-us/azure/role-based-access-control/conditions-overview#what-are-role-assignment-conditions) for details. Changing this forces a new resource to be created.

* `condition_version` - (Optional) The version of the condition. Supported values include `2.0`. Changing this forces a new resource to be created.

~> **Note:** `condition_version` is required when specifying `condition` and vice versa.

---

An `expiration` block supports the following:

* `duration_days` - (Optional) The duration of the role assignment in days. Changing this forces a new resource to be created.

* `duration_hours` - (Optional) The duration of the role assignment in hours. Changing this forces a new resource to be created.

* `end_date_time` - (Optional) The end date/time of the role assignment. Changing this forces a new resource to be created.

~> **Note:** Only one of `duration_days`, `duration_hours` or `end_date_time` should be specified.

---

A `schedule` block supports the following:

* `expiration` - (Optional) An `expiration` block as defined above.

* `start_date_time` - (Optional) The start date/time of the role assignment. Changing this forces a new resource to be created.

---

A `ticket` block supports the following:

* `number` - (Optional) User-supplied ticket number to be included with the request. Changing this forces a new resource to be created.

* `system` - (Optional) User-supplied ticket system name to be included with the request. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the PIM Eligible Role Assignment.
* `principal_type` - Type of principal to which the role will be assigned.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 minutes) Used when creating the PIM Eligible Role Assignment.
* `read` - (Defaults to 5 minutes) Used when retrieving the PIM Eligible Role Assignment.
* `delete` - (Defaults to 10 minutes) Used when deleting the PIM Eligible Role Assignment.

## Import

PIM Eligible Role Assignments can be imported using the following composite resource ID, e.g.

```shell
terraform import azurerm_pim_eligible_role_assignment.example /subscriptions/00000000-0000-0000-0000-000000000000|/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/roleDefinitions/00000000-0000-0000-0000-000000000000|00000000-0000-0000-0000-000000000000
```

-> **Note:** This ID is specific to Terraform - and is of the format `{scope}|{roleDefinitionId}|{principalId}`, where the first segment is the scope of the role assignment, the second segment is the role definition ID, and the last segment is the principal object ID.
