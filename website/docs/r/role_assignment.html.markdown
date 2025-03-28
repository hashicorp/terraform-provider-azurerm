---
subcategory: "Authorization"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_role_assignment"
description: |-
  Assigns a given Principal (User or Group) to a given Role.

---

# azurerm_role_assignment

Assigns a given Principal (User or Group) to a given Role.

## Example Usage (using a built-in Role)

```hcl
data "azurerm_subscription" "primary" {
}

data "azurerm_client_config" "example" {
}

resource "azurerm_role_assignment" "example" {
  scope                = data.azurerm_subscription.primary.id
  role_definition_name = "Reader"
  principal_id         = data.azurerm_client_config.example.object_id
}
```

## Example Usage (Custom Role & Service Principal)

```hcl
data "azurerm_subscription" "primary" {
}

data "azurerm_client_config" "example" {
}

resource "azurerm_role_definition" "example" {
  role_definition_id = "00000000-0000-0000-0000-000000000000"
  name               = "my-custom-role-definition"
  scope              = data.azurerm_subscription.primary.id

  permissions {
    actions     = ["Microsoft.Resources/subscriptions/resourceGroups/read"]
    not_actions = []
  }

  assignable_scopes = [
    data.azurerm_subscription.primary.id,
  ]
}

resource "azurerm_role_assignment" "example" {
  name               = "00000000-0000-0000-0000-000000000000"
  scope              = data.azurerm_subscription.primary.id
  role_definition_id = azurerm_role_definition.example.role_definition_resource_id
  principal_id       = data.azurerm_client_config.example.object_id
}
```

## Example Usage (Custom Role & User)

```hcl
data "azurerm_subscription" "primary" {
}

data "azurerm_client_config" "example" {
}

resource "azurerm_role_definition" "example" {
  role_definition_id = "00000000-0000-0000-0000-000000000000"
  name               = "my-custom-role-definition"
  scope              = data.azurerm_subscription.primary.id

  permissions {
    actions     = ["Microsoft.Resources/subscriptions/resourceGroups/read"]
    not_actions = []
  }

  assignable_scopes = [
    data.azurerm_subscription.primary.id,
  ]
}

resource "azurerm_role_assignment" "example" {
  name               = "00000000-0000-0000-0000-000000000000"
  scope              = data.azurerm_subscription.primary.id
  role_definition_id = azurerm_role_definition.example.role_definition_resource_id
  principal_id       = data.azurerm_client_config.example.object_id
}
```

## Example Usage (Custom Role & Management Group)

```hcl
data "azurerm_subscription" "primary" {
}

data "azurerm_client_config" "example" {
}

data "azurerm_management_group" "example" {
  name = "00000000-0000-0000-0000-000000000000"
}

resource "azurerm_role_definition" "example" {
  role_definition_id = "00000000-0000-0000-0000-000000000000"
  name               = "my-custom-role-definition"
  scope              = data.azurerm_subscription.primary.id

  permissions {
    actions     = ["Microsoft.Resources/subscriptions/resourceGroups/read"]
    not_actions = []
  }

  assignable_scopes = [
    data.azurerm_subscription.primary.id,
  ]
}

resource "azurerm_role_assignment" "example" {
  name               = "00000000-0000-0000-0000-000000000000"
  scope              = data.azurerm_management_group.primary.id
  role_definition_id = azurerm_role_definition.example.role_definition_resource_id
  principal_id       = data.azurerm_client_config.example.object_id
}
```

## Example Usage (ABAC Condition)

```hcl
data "azurerm_subscription" "primary" {
}

data "azurerm_client_config" "example" {
}

data "azurerm_role_definition" "builtin" {
  name = "Reader"
}

resource "azurerm_role_assignment" "example" {
  role_definition_name = "Role Based Access Control Administrator"
  scope                = data.azurerm_subscription.primary.id
  principal_id         = data.azurerm_client_config.example.object_id
  principal_type       = "ServicePrincipal"
  description          = "Role Based Access Control Administrator role assignment with ABAC Condition."
  condition_version    = "2.0"
  condition            = <<-EOT
(
 (
  !(ActionMatches{'Microsoft.Authorization/roleAssignments/write'})
 )
 OR
 (
  @Request[Microsoft.Authorization/roleAssignments:RoleDefinitionId] ForAnyOfAnyValues:GuidEquals {${basename(data.azurerm_role_definition.builtin.role_definition_id)}}
 )
)
AND
(
 (
  !(ActionMatches{'Microsoft.Authorization/roleAssignments/delete'})
 )
 OR
 (
  @Resource[Microsoft.Authorization/roleAssignments:RoleDefinitionId] ForAnyOfAnyValues:GuidEquals {${basename(data.azurerm_role_definition.builtin.role_definition_id)}}
 )
)
EOT
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) A unique UUID/GUID for this Role Assignment - one will be generated if not specified. Changing this forces a new resource to be created.

* `scope` - (Required) The scope at which the Role Assignment applies to, such as `/subscriptions/0b1f6471-1bf0-4dda-aec3-111122223333`, `/subscriptions/0b1f6471-1bf0-4dda-aec3-111122223333/resourceGroups/myGroup`, or `/subscriptions/0b1f6471-1bf0-4dda-aec3-111122223333/resourceGroups/myGroup/providers/Microsoft.Compute/virtualMachines/myVM`, or `/providers/Microsoft.Management/managementGroups/myMG`. Changing this forces a new resource to be created.

* `role_definition_id` - (Optional) The Scoped-ID of the Role Definition. Changing this forces a new resource to be created.

* `role_definition_name` - (Optional) The name of a built-in Role. Changing this forces a new resource to be created.

~> **Note:** Either `role_definition_id` or `role_definition_name` must be set.

* `principal_id` - (Required) The ID of the Principal (User, Group or Service Principal) to assign the Role Definition to. Changing this forces a new resource to be created.

-> **Note:** The Principal ID is also known as the Object ID (i.e. not the "Application ID" for applications).

* `principal_type` - (Optional) The type of the `principal_id`. Possible values are `User`, `Group` and `ServicePrincipal`. Changing this forces a new resource to be created. It is necessary to explicitly set this attribute when creating role assignments if the principal creating the assignment is constrained by ABAC rules that filters on the PrincipalType attribute.

* `condition` - (Optional) The condition that limits the resources that the role can be assigned to. Changing this forces a new resource to be created.

* `condition_version` - (Optional) The version of the condition. Possible values are `1.0` or `2.0`. Changing this forces a new resource to be created.

~> **Note:** `condition` is required when `condition_version` is set.

* `delegated_managed_identity_resource_id` - (Optional) The delegated Azure Resource Id which contains a Managed Identity. Changing this forces a new resource to be created.

~> **Note:** This field is only used in cross tenant scenarios.

* `description` - (Optional) The description for this Role Assignment. Changing this forces a new resource to be created.

* `skip_service_principal_aad_check` - (Optional) If the `principal_id` is a newly provisioned `Service Principal` set this value to `true` to skip the `Azure Active Directory` check which may fail due to replication lag. This argument is only valid if the `principal_id` is a `Service Principal` identity. Defaults to `false`.

~> **Note:** If it is not a `Service Principal` identity it will cause the role assignment to fail.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The Role Assignment ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Role Assignment.
* `read` - (Defaults to 5 minutes) Used when retrieving the Role Assignment.
* `delete` - (Defaults to 30 minutes) Used when deleting the Role Assignment.

## Import

Role Assignments can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_role_assignment.example /subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/roleAssignments/00000000-0000-0000-0000-000000000000
```

~> **Note:** The format of `resource id` could be different for different kinds of `scope`:

* for scope `Subscription`, the id format is `/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/roleAssignments/00000000-0000-0000-0000-000000000000`
* for scope `Resource Group`, the id format is `/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Authorization/roleAssignments/00000000-0000-0000-0000-000000000000`
* for scope `Key Vault`, the id format is `/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.KeyVault/vaults/vault1/providers/Microsoft.Authorization/roleAssignments/00000000-0000-0000-0000-000000000000`
* for scope `Storage Account`, the id format is `/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Storage/storageAccounts/storageAccount1/providers/Microsoft.Authorization/roleAssignments/00000000-0000-0000-0000-000000000000`

~> **Note:** for cross tenant scenarios, the format of `resource id` is composed of Azure resource ID and tenantId. for example:

```text
/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/roleAssignments/00000000-0000-0000-0000-000000000000|00000000-0000-0000-0000-000000000000
```
