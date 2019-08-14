---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_role_assignment"
sidebar_current: "docs-azurerm-resource-authorization-role-assignment"
description: |-
  Assigns a given Principal (User or Application) to a given Role.

---

# azurerm_role_assignment

Assigns a given Principal (User or Application) to a given Role.

## Example Usage (using a built-in Role)

```hcl
data "azurerm_subscription" "primary" {}

data "azurerm_client_config" "test" {}

resource "azurerm_role_assignment" "test" {
  scope                = "${data.azurerm_subscription.primary.id}"
  role_definition_name = "Reader"
  principal_id         = "${data.azurerm_client_config.test.service_principal_object_id}"
}
```

## Example Usage (Custom Role & Service Principal)

```hcl
data "azurerm_subscription" "primary" {}

data "azurerm_client_config" "test" {}

resource "azurerm_role_definition" "test" {
  role_definition_id = "00000000-0000-0000-0000-000000000000"
  name               = "my-custom-role-definition"
  scope              = "${data.azurerm_subscription.primary.id}"

  permissions {
    actions     = ["Microsoft.Resources/subscriptions/resourceGroups/read"]
    not_actions = []
  }

  assignable_scopes = [
    "${data.azurerm_subscription.primary.id}",
  ]
}

resource "azurerm_role_assignment" "test" {
  name               = "00000000-0000-0000-0000-000000000000"
  scope              = "${data.azurerm_subscription.primary.id}"
  role_definition_id = "${azurerm_role_definition.test.id}"
  principal_id       = "${data.azurerm_client_config.test.service_principal_object_id}"
}
```

## Example Usage (Custom Role & User)

```hcl
data "azurerm_subscription" "primary" {}

data "azurerm_client_config" "test" {}

resource "azurerm_role_definition" "test" {
  role_definition_id = "00000000-0000-0000-0000-000000000000"
  name               = "my-custom-role-definition"
  scope              = "${data.azurerm_subscription.primary.id}"

  permissions {
    actions     = ["Microsoft.Resources/subscriptions/resourceGroups/read"]
    not_actions = []
  }

  assignable_scopes = [
    "${data.azurerm_subscription.primary.id}",
  ]
}

resource "azurerm_role_assignment" "test" {
  name               = "00000000-0000-0000-0000-000000000000"
  scope              = "${data.azurerm_subscription.primary.id}"
  role_definition_id = "${azurerm_role_definition.test.id}"
  principal_id       = "${data.azurerm_client_config.test.client_id}"
}
```

## Example Usage (Custom Role & Management Group)

```hcl
data "azurerm_subscription" "primary" {}

data "azurerm_client_config" "test" {}

data "azurerm_management_group" "test" {}

resource "azurerm_role_definition" "test" {
  role_definition_id = "00000000-0000-0000-0000-000000000000"
  name               = "my-custom-role-definition"
  scope              = "${data.azurerm_subscription.primary.id}"

  permissions {
    actions     = ["Microsoft.Resources/subscriptions/resourceGroups/read"]
    not_actions = []
  }

  assignable_scopes = [
    "${data.azurerm_subscription.primary.id}",
  ]
}

resource "azurerm_role_assignment" "test" {
  name               = "00000000-0000-0000-0000-000000000000"
  scope              = "${data.azurerm_management_group.primary.id}"
  role_definition_id = "${azurerm_role_definition.test.id}"
  principal_id       = "${data.azurerm_client_config.test.client_id}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) A unique UUID/GUID for this Role Assignment - one will be generated if not specified. Changing this forces a new resource to be created.

* `scope` - (Required) The scope at which the Role Assignment applies to, such as `/subscriptions/0b1f6471-1bf0-4dda-aec3-111122223333`, `/subscriptions/0b1f6471-1bf0-4dda-aec3-111122223333/resourceGroups/myGroup`, or `/subscriptions/0b1f6471-1bf0-4dda-aec3-111122223333/resourceGroups/myGroup/providers/Microsoft.Compute/virtualMachines/myVM`, or `/providers/Microsoft.Management/managementGroups/myMG`. Changing this forces a new resource to be created.

* `role_definition_id` - (Optional) The Scoped-ID of the Role Definition. Changing this forces a new resource to be created. Conflicts with `role_definition_name`.

* `role_definition_name` - (Optional) The name of a built-in Role. Changing this forces a new resource to be created. Conflicts with `role_definition_id`.

* `principal_id` - (Required) The ID of the Principal (User, Group, Service Principal, or Application) to assign the Role Definition to. Changing this forces a new resource to be created. 

~> **NOTE:** The Principal ID is also known as the Object ID (ie not the "Application ID" for applications).

## Attributes Reference

The following attributes are exported:

* `id` - The Role Assignment ID.

## Import

Role Assignments can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_role_assignment.test /subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/roleAssignments/00000000-0000-0000-0000-000000000000
```
