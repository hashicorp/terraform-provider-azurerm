---
subcategory: "Authorization"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_role_definition"
description: |-
  Manages a custom Role Definition.

---

# azurerm_role_definition

Manages a custom Role Definition, used to assign Roles to Users/Principals. See ['Understand role definitions'](https://docs.microsoft.com/en-us/azure/role-based-access-control/role-definitions) in the Azure documentation for more details.

## Example Usage

```hcl
data "azurerm_subscription" "primary" {
}

resource "azurerm_role_definition" "example" {
  name        = "my-custom-role"
  scope       = data.azurerm_subscription.primary.id
  description = "This is a custom role created via Terraform"

  permissions {
    actions     = ["*"]
    not_actions = []
  }

  assignable_scopes = [
    data.azurerm_subscription.primary.id, # /subscriptions/00000000-0000-0000-0000-000000000000
  ]
}
```

## Argument Reference

The following arguments are supported:

* `role_definition_id` - (Optional) A unique UUID/GUID which identifies this role - one will be generated if not specified. Changing this forces a new resource to be created.

* `name` - (Required) The name of the Role Definition. Changing this forces a new resource to be created.

* `scope` - (Required) The scope at which the Role Definition applies too, such as `/subscriptions/0b1f6471-1bf0-4dda-aec3-111122223333`, `/subscriptions/0b1f6471-1bf0-4dda-aec3-111122223333/resourceGroups/myGroup`, or `/subscriptions/0b1f6471-1bf0-4dda-aec3-111122223333/resourceGroups/myGroup/providers/Microsoft.Compute/virtualMachines/myVM`. It is recommended to use the first entry of the `assignable_scopes`. Changing this forces a new resource to be created.

* `description` - (Optional) A description of the Role Definition.

* `permissions` - (Required) A `permissions` block as defined below.

* `assignable_scopes` - (Optional) One or more assignable scopes for this Role Definition, such as `/subscriptions/0b1f6471-1bf0-4dda-aec3-111122223333`, `/subscriptions/0b1f6471-1bf0-4dda-aec3-111122223333/resourceGroups/myGroup`, or `/subscriptions/0b1f6471-1bf0-4dda-aec3-111122223333/resourceGroups/myGroup/providers/Microsoft.Compute/virtualMachines/myVM`.

~> **NOTE:** The value for `scope` is automatically included in this list.

A `permissions` block as the following properties:

* `actions` - (Optional) One or more Allowed Actions, such as `*`, `Microsoft.Resources/subscriptions/resourceGroups/read`. See ['Azure Resource Manager resource provider operations'](https://docs.microsoft.com/en-us/azure/role-based-access-control/resource-provider-operations) for details.

* `data_actions` - (Optional) One or more Allowed Data Actions, such as `*`, `Microsoft.Storage/storageAccounts/blobServices/containers/blobs/read`. See ['Azure Resource Manager resource provider operations'](https://docs.microsoft.com/en-us/azure/role-based-access-control/resource-provider-operations) for details.

* `not_actions` - (Optional) One or more Disallowed Actions, such as `*`, `Microsoft.Resources/subscriptions/resourceGroups/read`. See ['Azure Resource Manager resource provider operations'](https://docs.microsoft.com/en-us/azure/role-based-access-control/resource-provider-operations) for details.

* `not_data_actions` - (Optional) One or more Disallowed Data Actions, such as `*`, `Microsoft.Resources/subscriptions/resourceGroups/read`. See ['Azure Resource Manager resource provider operations'](https://docs.microsoft.com/en-us/azure/role-based-access-control/resource-provider-operations) for details.

## Attributes Reference

The following attributes are exported:

* `id` - The Role Definition ID.

* `role_definition_id` - The Azure Resource Manager ID for the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Role Definition.
* `update` - (Defaults to 30 minutes) Used when updating the Role Definition.
* `read` - (Defaults to 5 minutes) Used when retrieving the Role Definition.
* `delete` - (Defaults to 30 minutes) Used when deleting the Role Definition.

## Import

Role Definitions can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_role_definition.example "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/roleDefinitions/00000000-0000-0000-0000-000000000000|/subscriptions/00000000-0000-0000-0000-000000000000"
```

-> **NOTE:** This ID is specific to Terraform - and is of the format `{roleDefinitionId}|{scope}`.
