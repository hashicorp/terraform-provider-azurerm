---
subcategory: "Managed Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_registration_definition"
description: |-
  Manages a Registration Definition.

---

# azurerm_registration_definition

Manages a Registration Definition.

## Example Usage

```hcl
data "azurerm_subscription" "primary" {
}

resource "azurerm_registration_definition" "example" {
  registration_definition_name  = "Sample registration"
  scope                         = data.azurerm_subscription.primary.id
  description                   = "This is a registration definition created via Terraform"
  managed_by_tenant_id          = "00000000-0000-0000-0000-000000000000"

  authorization {
    principal_id        = "00000000-0000-0000-0000-000000000000"
    role_definition_id  = "b24988ac-6180-42a0-ab88-20f7382dd24c"
  }
}
```

## Argument Reference

The following arguments are supported:

* `registration_definition_id` - (Optional) A unique UUID/GUID which identifies this registration definition - one will be generated if not specified. Changing this forces a new resource to be created.

* `registration_definition_name` - (Required) The name of the Registration Definition.

* `scope` - (Required) The scope at which the Registration Definition applies too, such as `/subscriptions/0b1f6471-1bf0-4dda-aec3-111122223333`. Only subscription level scope is supported.

* `managed_by_tenant_id` - (Required) A ID of the managing tenant.

* `description` - (Optional) A description of the Registration Definition.

* `authorization` - (Required) Authorization tuple containing principal id of the user/security group or service principal and id of the build-in role.

A `authorization` block as the following properties:

* `principal_id` - (Required) Principal Id of the security group/service principal/user that would be assigned permissions to the projected subscription or resource group.

* `role_definition_id` - (Required) The role definition identifier. This role will define all the permissions that the security group/service principal/user must have on the projected subscription or resource group. This role cannot be an owner role.

## Attributes Reference

The following attributes are exported:

* `id` - The Registration Definition ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Registration Definition.
* `update` - (Defaults to 30 minutes) Used when updating the Registration Definition.
* `read` - (Defaults to 5 minutes) Used when retrieving the Registration Definition.
* `delete` - (Defaults to 30 minutes) Used when deleting the Registration Definition.

## Import

Registration Definitions can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_registration_definition.example /subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.ManagedServices/registrationDefinitions/00000000-0000-0000-0000-000000000000
```
