---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_role_assignment"
sidebar_current: "docs-azurerm-resource-authorization-role-assignment"
description: |-
  Assigns a given Principal (User or Application) to a given Role.

---

# azurerm_role_assignment

Assigns a given Principal (User or Application) to a given Role.

## Example Usage

Complete examples of how to use the `azurerm_role_assignment` resource can be found [in the `./examples/roles` folder within the Github Repository](https://github.com/terraform-providers/terraform-provider-azurerm/tree/master/examples/roles)

```hcl
data "azurerm_subscription" "example" {}

data "azurerm_client_config" "example" {}

resource "azurerm_role_assignment" "example" {
  scope                = "${data.azurerm_subscription.example.id}"
  role_definition_name = "Reader"
  principal_id         = "${data.azurerm_client_config.example.service_principal_object_id}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) A unique UUID/GUID for this Role Assignment - one will be generated if not specified. Changing this forces a new resource to be created.

* `scope` - (Required) The scope at which the Role Assignment applies to. This can be the ID of a Subscription (e.g. `/subscriptions/0b1f6471-1bf0-4dda-aec3-111122223333`), a Resource Group (e.g. `/subscriptions/0b1f6471-1bf0-4dda-aec3-111122223333/resourceGroups/myGroup`) or a resource within a Resource Group (e.g. `/subscriptions/0b1f6471-1bf0-4dda-aec3-111122223333/resourceGroups/myGroup/providers/Microsoft.Compute/virtualMachines/myVM`). Changing this forces a new resource to be created.

-> **NOTE:** We recommend using [Interpolation Syntax](https://www.terraform.io/docs/configuration/interpolation.html) to pull this value from a Data Source or Resource where possible, rather than hard-coding the ID's.

* `role_definition_id` - (Optional) The Scoped-ID of the Role Definition. Changing this forces a new resource to be created. Conflicts with `role_definition_name`.

* `role_definition_name` - (Optional) The name of a built-in Role. Changing this forces a new resource to be created. Conflicts with `role_definition_id`.

* `principal_id` - (Required) The ID of the Principal (User or Application) to assign the Role Definition to. Changing this forces a new resource to be created.


## Attributes Reference

The following attributes are exported:

* `id` - The Role Assignment ID.

## Import

Role Assignments can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_role_assignment.test /subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/roleAssignments/00000000-0000-0000-0000-000000000000
```
