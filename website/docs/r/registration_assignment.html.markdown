---
subcategory: "Managed Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_registration_assignment"
description: |-
    Assigns a given Registration Definition to a subscription or a resource group.

---

# azurerm_registration_assignment

Assigns a given Registration Definition to a subscription or a resource group.

## Example Usage

```hcl
data "azurerm_subscription" "primary" {
}

resource "azurerm_registration_definition" "example" {
  registration_definition_name = "Sample registration"
  description                  = "This is a registration definition created via Terraform"
  managed_by_tenant_id         = "00000000-0000-0000-0000-000000000000"

  authorization {
    principal_id       = "00000000-0000-0000-0000-000000000000"
    role_definition_id = "b24988ac-6180-42a0-ab88-20f7382dd24c"
  }
}

resource "azurerm_registration_assignment" "example" {
  registration_assignment_id = "%s"
  scope                      = data.azurerm_subscription.primary.id
  registration_definition_id = azurerm_registration_definition.example.id
}
```

## Argument Reference

The following arguments are supported:

* `registration_assignment_id` - (Optional) A unique UUID/GUID which identifies this registration assignment- one will be generated if not specified. Changing this forces a new resource to be created.

* `scope` - (Required) The scope at which the Registration Assignment applies too, such as `/subscriptions/0b1f6471-1bf0-4dda-aec3-111122223333` or `/subscriptions/0b1f6471-1bf0-4dda-aec3-111122223333/resourceGroups/myGroup`. Changing this forces a new resource to be created.

* `registration_definition_id` - (Required) A Fully qualified path of the registration definitio, such as `/subscriptions/0afefe50-734e-4610-8c82-a144aff49dea/providers/Microsoft.ManagedServices/registrationDefinitions/26c128c2-fefa-4340-9bb1-8e081c90ada2`. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The Registration Assignment ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Registration Assignment.
* `update` - (Defaults to 30 minutes) Used when updating the Registration Assignment.
* `read` - (Defaults to 5 minutes) Used when retrieving the Registration Assignment.
* `delete` - (Defaults to 30 minutes) Used when deleting the Registration Assignment.

## Import

Registration Assignments can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_registration_assignment.example /subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.ManagedServices/registrationAssignments/00000000-0000-0000-0000-000000000000
```
