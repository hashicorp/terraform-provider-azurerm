---
subcategory: "Authorization"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_pim_eligible_role_assignment"
description: |-
  Manages a Pim Eligible Role Assignment.
---

# azurerm_pim_eligible_role_assignment

Manages a Pim Eligible Role Assignment.

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

* `principal_id` - (Required) The principal id. Changing this forces a new Pim Eligible Role Assignment to be created.

* `principal_type` - (Required) The type of principal. Changing this forces a new Pim Eligible Role Assignment to be created.

* `role_definition_id` - (Required) The role definition id. Changing this forces a new Pim Eligible Role Assignment to be created.

* `scope` - (Required) The scope. Changing this forces a new Pim Eligible Role Assignment to be created.

---

* `justification` - (Optional) The justification of the role assignment. Changing this forces a new Pim Eligible Role Assignment to be created.

* `schedule` - (Optional) A `schedule` block as defined below. Changing this forces a new Pim Eligible Role Assignment to be created.

* `ticket` - (Optional) A `ticket` block as defined below. Changing this forces a new Pim Eligible Role Assignment to be created.

---


A `expiration` block supports the following:

* `duration_days` - (Optional) The duration of the role assignment in days. Conflicts with `schedule.0.expiration.0.duration_hours`,`schedule.0.expiration.0.end_date_time` Changing this forces a new Pim Eligible Role Assignment to be created.

* `duration_hours` - (Optional) The duration of the role assignment in hours. Conflicts with `schedule.0.expiration.0.duration_days`,`schedule.0.expiration.0.end_date_time` Changing this forces a new Pim Eligible Role Assignment to be created.

* `end_date_time` - (Optional) The end date time of the role assignment. Conflicts with `schedule.0.expiration.0.duration_days`,`schedule.0.expiration.0.duration_hours` Changing this forces a new Pim Eligible Role Assignment to be created.

---

A `schedule` block supports the following:

* `expiration` - (Optional) A `expiration` block as defined above.

* `start_date_time` - (Optional) The start date time of the role assignment. Changing this forces a new Pim Eligible Role Assignment to be created.

---

A `ticket` block supports the following:

* `number` - (Optional) The ticket number.

* `system` - (Optional) The ticket system.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Pim Eligible Role Assignment.
* `principal_type` - The type of principal.
*
## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Pim Eligible Role Assignment.
* `read` - (Defaults to 5 minutes) Used when retrieving the Pim Eligible Role Assignment.
* `delete` - (Defaults to 30 minutes) Used when deleting the Pim Eligible Role Assignment.

## Import

Pim Eligible Role Assignments can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_pim_eligible_role_assignment.example /subscriptions/00000000-0000-0000-0000-000000000000|/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/roleDefinitions/00000000-0000-0000-0000-000000000000|00000000-0000-0000-0000-000000000000
```

-> **NOTE:** This ID is specific to Terraform - and is of the format `{scope}|{roleDefinitionId}|{principalId}`.
