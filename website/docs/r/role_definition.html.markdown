---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_role_definition"
sidebar_current: "docs-azurerm-resource-authorization-role-definition"
description: |-
  Manages a custom Role Definition.

---

# azurerm_role_definition

Manages a custom Role Definition, used to assign Roles to Users/Principals.

## Example Usage

```hcl
data "azurerm_subscription" "example" {}

resource "azurerm_role_definition" "example" {
  name               = "example-role"
  scope              = "${data.azurerm_subscription.example.id}"
  description        = "This is a custom role created via Terraform"

  permissions {
    actions     = ["*"]
    not_actions = []
  }

  assignable_scopes = [
    # /subscriptions/00000000-0000-0000-0000-000000000000
    "${data.azurerm_subscription.example.id}",
  ]
}
```

## Argument Reference

The following arguments are supported:

* `role_definition_id` - (Optional) A unique UUID/GUID which identifies this role - one will be generated if not specified. Changing this forces a new resource to be created.

* `name` - (Required) The name of the Role Definition. Changing this forces a new resource to be created.

* `scope` - (Required) The scope at which the Role Definition applies to. This can be the ID of a Subscription (e.g. `/subscriptions/0b1f6471-1bf0-4dda-aec3-111122223333`), a Resource Group (e.g. `/subscriptions/0b1f6471-1bf0-4dda-aec3-111122223333/resourceGroups/myGroup`) or a resource within a Resource Group (e.g. `/subscriptions/0b1f6471-1bf0-4dda-aec3-111122223333/resourceGroups/myGroup/providers/Microsoft.Compute/virtualMachines/myVM`). Changing this forces a new resource to be created.

-> **NOTE:** We recommend using [Interpolation Syntax](https://www.terraform.io/docs/configuration/interpolation.html) to pull this value from a Data Source or Resource where possible, rather than hard-coding the ID's.


* `description` - (Optional) A description of the Role Definition.

* `permissions` - (Required) A `permissions` block as defined below.

* `assignable_scopes` - (Required) One or more assignable scopes for this Role Definition. This can be the ID of a Subscription (e.g. `/subscriptions/0b1f6471-1bf0-4dda-aec3-111122223333`), a Resource Group (e.g. `/subscriptions/0b1f6471-1bf0-4dda-aec3-111122223333/resourceGroups/myGroup`) or a resource within a Resource Group (e.g. `/subscriptions/0b1f6471-1bf0-4dda-aec3-111122223333/resourceGroups/myGroup/providers/Microsoft.Compute/virtualMachines/myVM`). Changing this forces a new resource to be created.

-> **NOTE:** We recommend using [Interpolation Syntax](https://www.terraform.io/docs/configuration/interpolation.html) to pull this value from a Data Source or Resource where possible, rather than hard-coding the ID's.

A `permissions` block as the following properties:

* `action` - (Optional) One or more Allowed Actions, such as `*`, `Microsoft.Resources/subscriptions/resourceGroups/read`.

* `not_action` - (Optional) One or more Disallowed Actions, such as `*`, `Microsoft.Resources/subscriptions/resourceGroups/read`.

## Attributes Reference

The following attributes are exported:

* `id` - The Role Definition ID.

## Import

Role Definitions can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_role_definition.test /subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/roleDefinitions/00000000-0000-0000-0000-000000000000
```
