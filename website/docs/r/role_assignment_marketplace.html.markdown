---
subcategory: "Authorization"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_role_assignment_marketplace"
description: |-
  Assigns a given Principal (User or Group) to a given Role at marketplace scope.
---

# azurerm_role_assignment_marketplace

Assigns a given Principal (User or Group) to a given Role at marketplace scope.

## Example Usage

```hcl
data "azurerm_client_config" "example" {
}

resource "azurerm_role_assignment_marketplace" "example" {
  role_definition_name = "Marketplace Admin"
  principal_id         = data.azurerm_client_config.example.object_id
}
```

## Argument Reference

The following arguments are supported:

* `principal_id` - (Required) The ID of the Principal (User, Group or Service Principal) to assign the Role Definition to. Changing this forces a new resource to be created.

* `name` - (Optional) A unique UUID/GUID for this Role Assignment - one will be generated if not specified. Changing this forces a new resource to be created.

* `role_definition_id` - (Optional) The Scoped-ID of the Role Definition. Changing this forces a new resource to be created. Conflicts with `role_definition_name`.

* `role_definition_name` - (Optional) The name of a built-in Role. Changing this forces a new resource to be created. Conflicts with `role_definition_id`.

~> **NOTE:** The Principal ID is also known as the Object ID (i.e. not the "Application ID" for applications). To assign Azure roles, the Principal must have `Microsoft.Authorization/roleAssignments/write` permissions. See [documentation](https://learn.microsoft.com/en-us/azure/role-based-access-control/role-assignments-portal) for more information. The calling principal must first be assigned Privileged Role Administrator (like `Owner` role) or Global Administrator. See [documentation](https://learn.microsoft.com/en-us/marketplace/create-manage-private-azure-marketplace-new#prerequisites) for more information.

* `delegated_managed_identity_resource_id` - (Optional) The delegated Azure Resource ID which contains a Managed Identity. Changing this forces a new resource to be created.

~> **NOTE:** this field is only used in cross tenant scenario.

* `description` - (Optional) The description for this Role Assignment. Changing this forces a new resource to be created.
  
* `skip_service_principal_aad_check` - (Optional) If the `principal_id` is a newly provisioned `Service Principal` set this value to `true` to skip the `Azure Active Directory` check which may fail due to replication lag. This argument is only valid if the `principal_id` is a `Service Principal` identity. Defaults to `false`.

~> **NOTE:** If it is not a `Service Principal` identity it will cause the role assignment to fail.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The Role Assignment ID.

* `principal_type` - The type of the `principal_id`, e.g. User, Group, Service Principal, Application, etc.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Role Assignment.
* `update` - (Defaults to 30 minutes) Used when updating the Role Assignment.
* `read` - (Defaults to 5 minutes) Used when retrieving the Role Assignment.
* `delete` - (Defaults to 30 minutes) Used when deleting the Role Assignment.

## Import

Role Assignments can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_role_assignment_marketplace.example /providers/Microsoft.Marketplace/providers/Microsoft.Authorization/roleAssignments/00000000-0000-0000-0000-000000000000
```

